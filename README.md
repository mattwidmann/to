to
==

`to` is a command-line utility for managing a to do list.  By default, the file's name is `to.txt` and it's stored in `~/Documents/`.  To do list items are separated by newlines and are sorted alphabetically.

Without any arguments, `to` displays the items in a to do list along with their line numbers (used to remove the items from a list).

With unrecognized arguments, `to` creates a new to do item in the default list by concatenating its arguments into one item, joined by spaces.

The `-r` argument followed by a valid integer will remove the to do item with that number from the to do list.

The `-d` and `-n` arguments take a path to a directory and filename, respectively.  These change which to do list file the other actions modify or display.

Example
-------

	$ # to do lists start out empty
	$ ./to
	$ ./to This is a new to do list item.
	$ ./to
	> 0 - This is a new to do list item.
	$ ./to This is another to do list item.
	$ ./to
	> 0 - This is a new to do list item.
	> 1 - This is another to do list item.
	$ ./to -r 0
	$ ./to
	> 0 - This is a new to do list item.

Caveats
-------

Line numbers will shift downwards or upwards when an item is removed or added.  Line numbers, therefore, are relative identifiers and will change if the list is modified.
