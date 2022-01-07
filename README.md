# torrent-metainfo-parser
generates a .torrent meta info from a file <br/>

**required arguments :**  <br />
`filePath`
_flag arguments (optional) :_ <br />
`-p piece size (in KB)` <br />
`-c comment` <br />
`-b creator name` <br />
`-a announce Url` <br />
if not specified they will be replaced by some defined default values

example : <br />
`examples/iceberg.jpg -a 127.0.0.1:6969 -b Asaad -p 16`
