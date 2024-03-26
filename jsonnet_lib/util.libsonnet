{
  // Case   Result
  // /       /
  // /asdf   /
  // .       .
  // ..      .
  // ./      .
  // ../     .
  // foo/    .
  dirname(path):
    if path == '' then '.'
    else if path == '/' then '/'
    else if path == '.' || path == './' then '.'
    else if path == '..' || path == '../' then '.'
    else {
      local stripTrailingSlash = if path[std.length(path)-1] != '/' then path
                                 else path[0:std.length(path)-1],
      local splitOnSlash = std.split(stripTrailingSlash, '/'),
      local removeLast = if std.length(splitOnSlash) < 2 then ['.']
                         else splitOnSlash[0:std.length(splitOnSlash)-1],
      ret: if removeLast == [''] then '/' else std.join('/', removeLast),
    }.ret,
}
