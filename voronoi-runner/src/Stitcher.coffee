{spawn} = require 'child_process'
fs = require 'fs'

class module.exports
  constructor: (@opts) ->

  stitch: (cb) ->
    args = [
      '-r', @opts.rIn
      '-i', @opts.inPattern
      '-vcodec', 'libx264'
      '-vpre', 'lossless_ultrafast'
      '-r', @opts.rOut
      '-pix_fmt', 'yuv420p'
      @opts.out
    ]
    fs.unlink @opts.out, (err) ->
      # Ignore error.
      program = spawn 'ffmpeg', args
      program.stdout.on 'data', (data) -> process.stdout.write data
      program.stderr.on 'data', (data) -> process.stderr.write data
      program.on 'close', (code) ->
        return cb 'err-' + code unless code is 0
        cb()
