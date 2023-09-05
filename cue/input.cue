package main

#Instructions: {
	from?: =~ "cgr.dev/chainguard"
	env?: [...string]
	run?: [...string]
	workdir?: string
	copy?: [...string]
	user?: string
	cmd?: [...string]
	...
}

#Dockerfile: [...{
	stage: int
	instructions: [...#Instructions]
}]

dockerfile: #Dockerfile
