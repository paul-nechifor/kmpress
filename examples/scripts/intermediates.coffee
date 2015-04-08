fs = require 'fs'
{spawn} = require 'child_process'
{Runner} = require '../../runner'
{Voronoi, Stitcher} = require '../../voronoi-runner'

voronoiDir = __dirname + '/../results/voronoi-frames'

getIntermediaryValues = (image, cb) ->
  r = new Runner
  args = ['-i', image, '-o', '', '-max', '1000']
  r.run args, (err) ->
    return cb err if err
    grouped = (groupArray x for x in r.iterationColors)
    cb null, grouped

groupArray = (array) ->
  ret = []
  elem = []
  for e, i in array
    elem.push ((e/255) * 10.0) - 5.0
    if i % 3 is 2
      ret.push elem
      elem = []
  return ret

renderFrame = (array, i, cb) ->
  v = new Voronoi
  path = "#{voronoiDir}/#{i}.png"
  opts =
    width: 300
    height: 200
    highQuality: true
    camera:
      location: [30, -50, 25]
      look_at: [0, 0, -0.5]
      right: '-0.24*x*image_width/image_height'
      up: '0.24*z'

  v.run array, path, opts, (err) ->
    return cb err if err
    renderClusters array, i, cb

renderClusters = (array, i, cb) ->
  ns = array.map (n) -> n
  ns = []
  ns.push.apply ns, n for n in array
  ns = ns.map (n) -> Math.floor (n + 5) / 10 * 255
  ns.splice 0, 0, ns.length
  list = ns.join ' '

  voronoi = __dirname + '/../../kmpress/kmpress'
  image = __dirname + '/../images/bird_large.tiff'
  out = "#{voronoiDir}/clusters-#{i}.tiff"
  s = spawn voronoi, ['-i', image, '-o', out, '-renderCluster']
  s.stdin.end list
  s.on 'close', (code) ->
    return cb 'err-' + code unless code is 0
    cb()

renderAllSteps = (arrays, cb) ->
  i = 0
  next = ->
    return cb null if i >= arrays.length
    renderFrame arrays[i], i, (err) ->
      console.log 'Completed frame', i, 'of', arrays.length
      return cb err if err
      i++
      next()
  next()

interpolateEvolution = (clustersEvolution, nStages) ->
  ret = []

  prev = clustersEvolution[0]
  ret.push prev

  for i in [1 .. clustersEvolution.length - 1]
    clus = clustersEvolution[i]
    addInterpolated ret, prev, clus, nStages
    prev = clus
    ret.push prev

  return ret

# nStages represends the number of intermediary points.
addInterpolated = (list, start, end, nStages) ->
  n = start.length
  for i in [1 .. nStages]
    stage =
      for j in [0 .. n - 1]
        for k in [0 .. 2]
          increment = (end[j][k] - start[j][k]) / (1 + nStages)
          start[j][k] + i * increment
    list.push stage
  return

main = ->
  image = __dirname + '/../images/bird_large.tiff'
  getIntermediaryValues image, (err, clustersEvolution) ->
    throw err if err
    fs.mkdir voronoiDir, (err) ->
      # Ignore error so far.
      interpolated = interpolateEvolution clustersEvolution, 6
      renderAllSteps interpolated, (err) ->
        s = new Stitcher
          inPattern: voronoiDir + '/%d.png'
          rIn: 30
          rOut: 30
          out: __dirname + '/../results/voronoi.mp4'
        s.stitch (err) ->
          throw err if err

main()
