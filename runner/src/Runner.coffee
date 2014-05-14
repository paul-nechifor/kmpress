{spawn} = require 'child_process'
byline = require 'byline'

class module.exports
  constructor: ->
    @lineInterpreters =
      iteration: @onIteration.bind this
    @iterationColors = []

  run: (args, cb) ->
    program = spawn __dirname + '/../../kmpress', args
    lineStream = byline.createStream program.stdout
    lineStream.on 'data', @interpretLine.bind this
    program.on 'close', (code) ->
      return cb 'err-' + code unless code is 0
      cb()

  interpretLine: (line) ->
    line = line.toString()
    start = line.indexOf '>>>'
    return if start is -1
    type = line.substring 0, start
    data = line.substring start + 3, line.length
    @lineInterpreters[type]? data

  onIteration: (line) ->
    numbers = line.split('[')[1].split(']')[0].split(' ')
    @iterationColors.push (Number(x) for x in numbers)
