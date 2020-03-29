# Taskmaster

## Commands
add    clear  fg        open  quit    remove  restart   start   stop  update
avail  exit   maintail  pid   reload  reread  shutdown  status  tail  version



Esc[m	Turn off character attributes	SGR0
Esc[0m	Turn off character attributes	SGR0
Esc[1m	Turn bold mode on	SGR1
Esc[2m	Turn low intensity mode on	SGR2
Esc[4m	Turn underline mode on	SGR4
Esc[5m	Turn blinking mode on	SGR5
Esc[7m	Turn reverse video on	SGR7
Esc[8m	Turn invisible text mode on	SGR8

Esc[Line;Liner	Set top and bottom lines of a window	DECSTBM

Esc[ValueA	Move cursor up n lines	CUU
Esc[ValueB	Move cursor down n lines	CUD
Esc[ValueC	Move cursor right n lines	CUF
Esc[ValueD	Move cursor left n lines	CUB
Esc[H	Move cursor to upper left corner	cursorhome
Esc[;H	Move cursor to upper left corner	cursorhome
Esc[Line;ColumnH	Move cursor to screen location v,h	CUP
Esc[f	Move cursor to upper left corner	hvhome
Esc[;f	Move cursor to upper left corner	hvhome
Esc[Line;Columnf	Move cursor to screen location v,h	CUP
EscD	Move/scroll window up one line	IND
EscM	Move/scroll window down one line	RI
EscE	Move to next line	NEL
Esc7	Save cursor position and attributes	DECSC
Esc8	Restore cursor position and attributes	DECSC

EscH	Set a tab at the current column	HTS
Esc[g	Clear a tab at the current column	TBC
Esc[0g	Clear a tab at the current column	TBC
Esc[3g	Clear all tabs	TBC

Esc#3	Double-height letters, top half	DECDHL
Esc#4	Double-height letters, bottom half	DECDHL
Esc#5	Single width, single height letters	DECSWL
Esc#6	Double width, single height letters	DECDWL

Esc[K	Clear line from cursor right	EL0
Esc[0K	Clear line from cursor right	EL0
Esc[1K	Clear line from cursor left	EL1
Esc[2K	Clear entire line	EL2

Esc[J	Clear screen from cursor down	ED0
Esc[0J	Clear screen from cursor down	ED0
Esc[1J	Clear screen from cursor up	ED1
Esc[2J	Clear entire screen	ED2

Esc5n	Device status report	DSR
Esc0n	Response: terminal is OK	DSR
Esc3n	Response: terminal is not OK	DSR

Esc6n	Get cursor position	DSR
EscLine;ColumnR	Response: cursor is at v,h	CPR

Esc[c	Identify what terminal type	DA
Esc[0c	Identify what terminal type (another)	DA
Esc[?1;Value0c	Response: terminal type code n	DA

Escc	Reset terminal to initial state	RIS

Esc#8	Screen alignment display	DECALN
Esc[2;1y	Confidence power up test	DECTST
Esc[2;2y	Confidence loopback test	DECTST
Esc[2;9y	Repeat power up test	DECTST
Esc[2;10y	Repeat loopback test	DECTST

Esc[0q	Turn off all four leds	DECLL0
Esc[1q	Turn on LED #1	DECLL1
Esc[2q	Turn on LED #2	DECLL2
Esc[3q	Turn on LED #3	DECLL3
Esc[4q	Turn on LED #4	DECLL4
