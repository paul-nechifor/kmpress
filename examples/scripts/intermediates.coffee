{Runner} = require '../../runner'

r = new Runner
image = __dirname + '/../images/bird_large.tiff'
args = ['-i', image, '-o', '/tmp/a.tiff', '-max', '100']
r.run args, (err) ->
  throw err if err
