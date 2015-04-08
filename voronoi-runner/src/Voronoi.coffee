{spawn} = require 'child_process'
fs = require 'fs'
tmp = require 'tmp'

class module.exports
  voronoi = __dirname + '/../../voronoi'

  constructor: ->

  run: (points, imagePath, opts, cb) ->
    tmp.dir {unsafeCleanup: true}, (err, path) =>
      return cb err if err
      @writeFiles path, points, opts, (err, sceneFile) =>
        return cb err if err
        args = @parseOpts opts
        args.push '+O' + imagePath
        args.push sceneFile
        @runPovRay args, (err) ->
          return cb err if err
          cb()

  writeFiles: (path, points, opts, cb) ->
    list = path + '/list'
    scene = path + '/scene.pov'
    particles = path + '/particles.pov'
    cells = path + '/cells.pov'
    @writeList list, points, (err) =>
      return cb err if err
      @writeScene scene, particles, cells, opts, (err) =>
        return cb err if err
        @runVoronoi [list, cells, particles], (err) ->
          return cb err if err
          cb null, scene

  runVoronoi: (args, cb) ->
    programPath = voronoi + '/voronoi'
    program = spawn programPath, args
    program.on 'close', (code) ->
      return cb 'err-' + code unless code is 0
      cb()

  runPovRay: (args, cb) ->
    program = spawn 'povray', args
    program.stderr.on 'data', (data) ->
    program.on 'close', (code) ->
      return cb 'err-' + code unless code is 0
      cb()

  parseOpts: (opts) ->
    args = []
    if opts.highQuality
      args.push ['+A0.0001', '+R9', '-J']
    args.push '+W' + opts.width
    args.push '+H' + opts.height
    return args

  pointsToFile: (points) ->
    list = []
    for p, i in points
      list.push i + ' ' + p.join ' '
    list.push ''
    return list.join '\n'

  writeList: (path, points, cb) ->
    fs.writeFile path, @pointsToFile(points), (err) ->
      return cb err if err
      cb()

  writeScene: (path, particlesFile, cellsFile, opts, cb) ->
    data = """
      #version 3.6;

      // Right-handed coordinate system in which the z-axis points upwards
      camera {
        location <#{opts.camera.location.join ','}>
        sky z
        right #{opts.camera.right}
        up #{opts.camera.up}
        look_at <#{opts.camera.look_at.join ','}>
      }

      // White background
      background{rgb 1}

      // Two lights with slightly different colors
      light_source{<-8,-20,30> color rgb <0.77,0.75,0.75>}
      light_source{<25,-12,12> color rgb <0.38,0.40,0.40>}

      // Radius of the Voronoi cell network
      #declare r=0.05;

      // Radius of the particles
      #declare s=0.2;

      // Particles
      union{
        #include "#{particlesFile}"
        pigment{rgb <1,0.4,0.45>} finish{reflection 0.1 specular 0.3 ambient 0.42}
      }

      // Voronoi cells
      union{
        #include "#{cellsFile}"
        pigment{rgb 0.95} finish{specular 0.5 ambient 0.42}
      }
    """

    fs.writeFile path, data, (err) ->
      return cb err if err
      cb()
