package pgutil

func Quote(s string) string {
	buf := make([]byte, 0, 3*len(s)/2)
	buf = append(buf, 'E')
	buf = append(buf, '\'')
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		switch r {
		case '\'':
			buf = append(buf, `'`...)
			buf = append(buf, byte(r))
		case '\b':
			buf = append(buf, `\b`...)
		case '\f':
			buf = append(buf, `\f`...)
		case '\n':
			buf = append(buf, `\n`...)
		case '\r':
			buf = append(buf, `\r`...)
		case '\t':
			buf = append(buf, `\t`...)
		case '\\':
			buf = append(buf, `\\`...)
		default:
			buf = append(buf, byte(r))
		}
	}
	buf = append(buf, '\'')
	return string(buf)
}
