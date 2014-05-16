fs = require 'fs'
{Runner} = require '../../runner'
{Voronoi} = require '../../voronoi-runner'

voronoiDir = __dirname + '/../results/voronoi-frames'

getIntermediaryValues = (image, cb) ->
  r = new Runner
  args = ['-i', image, '-o', '', '-max', '1000']
  r.run args, (err) ->
    return cb err if err
    cb null, r.iterationColors

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
    width: 1080
    height: 1080
    highQuality: true
  v.run groupArray(array), path, opts, cb

renderAllSteps = (arrays, cb) ->
  i = 0
  next = ->
    return cb null if i >= arrays.length
    renderFrame arrays[i], i, (err) ->
      console.log 'Completed frame:', i
      return cb err if err
      i++
      next()
  next()

main = ->
  image = __dirname + '/../images/bird_large.tiff'
  getIntermediaryValues image, (err, arrays) ->
    throw err if err
    fs.mkdir voronoiDir, (err) ->
      # Ignore error so far.
      renderAllSteps arrays, (err) ->
        throw err if err
      
main()
