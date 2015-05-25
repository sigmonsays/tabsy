package main

/*
 simple command line layer for making tab complete interactive CLIs

 cd [path]
 pwd

 show ?
 show object [path]

 delete [path]
 makedir [path]
 upload [localpath] [path]

 object set mtime [mtime]



*/

type ExecuteCommand func() error

type Command struct {
	Name string
	Exec ExecuteCommand

	SubCmd []Command
}
