This code was thrown together to get a handle on how to structure a project with 
multiple packages in a heirarchical structure.

A go.work file is created in the project to ensure file dependencies are not shown
as unavailable via the GOPATH.

The 'somemath' module highlights how multiple files can be used to define a single
logical package for referencing from dependant code. Note the 'special' subfolder
which uses it's own go.mod file to define the somemath/special path, and without 
which I was unable to make the 'special' file accessible to the calling main.go
file.

The go.mod files were required to make everything come together in terms of ensuring
paths to the source files would be clearly defined.  I'm guessing that the use of the 
'work' file requires the go mod files to be added to each folder to identify the 
include path names.  It feels a little clunky, but if this is the way it is supposed
to work then so be it.

In terms of the heirarchy defind here and the choices made to orgaise the code, 
I would not necessarily structure all of my projects this way. I can see that
in certain cases it might be beneficial to define such a heirarchy to help better 
organise/categorise functionality and define an api structure. However, it is likely 
to be overkill for most small projects, and was really only done here as a learning
exersize in order to better understand how modules and packages work in Golang.
