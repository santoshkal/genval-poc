package main

#Dockerfile: [{
	stage: int
	instructions: [{from: string},
				{env: [...string]},

		{run: [...string]},
		{workdir: string},
		{copy: [...string]},
		{run: [...string]},

		{cmd: [...string]},
		...]
}]
