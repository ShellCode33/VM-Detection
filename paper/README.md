# Paper

To build the paper, be sure to have Docker installed and run the following
command inside the paper directory:

```bash
$ docker run --rm -it -v "$(pwd):/pandoc" dalibo/pandocker --pdf-engine=xelatex --template=eisvogel --listings --highlight-style espresso *.md -o paper.pdf
```
