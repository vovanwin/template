package swagger

import (
	"html/template"
	"net/http"
)

type page struct {
	Path        string
	Description string
}

type indexPage struct {
	pages []page
}

func newIndexPage() *indexPage {
	return &indexPage{}
}

func (i *indexPage) addPage(path string, description string) {
	i.pages = append(i.pages, page{
		Path:        path,
		Description: description,
	})
}

func (i indexPage) handler(w http.ResponseWriter, list map[string]string, param Param) error {
	return template.Must(template.New("index").Parse(`<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>Elements in Twitter Bootstrap</title>

    <!-- Elements: Web Component -->
    <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">

    <!-- Twitter Bootstrap -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
    <!-- Twitter Bootstrap: Sticky Footer Example -->
    <link rel="stylesheet" href="https://getbootstrap.com/docs/4.5/examples/sticky-footer-navbar/sticky-footer-navbar.css">

    <style>
        body {
            display: flex;
            flex-direction: column;
            height: 100vh;
        }

        main {
            flex: 1 0 0;
            overflow: hidden;
        }
    </style>
</head>

<body>

<header>
    <!-- Fixed navbar -->
    <nav class="navbar navbar-expand-md navbar-dark bg-dark">
        <a class="navbar-brand" href="#">Api - Умный город</a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarCollapse">
            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <div class="dropdown">
                        <a class="btn btn-secondary dropdown-toggle" href="#" role="button" id="dropdownMenuLink" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            Выбери раздел документации
                        </a>

                        <div class="dropdown-menu" aria-labelledby="dropdownMenuLink">
                            {{range $index, $element :=  .List}}
			                <a class="dropdown-item" href="{{$element}}">{{$index}}</a>
			                {{end}}
                        </div>
                    </div>
                </li>
            </ul>
        </div>
    </nav>
</header>

<!-- Begin page content -->
<main role="main">
    <elements-api
            apiDescriptionUrl="{{.Path}}"
            router="memory"
            layout="sidebar"
    />

</main>

 
<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ho+j7jyWK8fNQe+A12Hb8AhRq26LrZ/JpcUGGOn+Y7RsweNrtN/tE3MoK7ZeZDyx" crossorigin="anonymous"></script>
</body>
</html>
`)).Execute(w, struct {
		Path string
		List map[string]string
	}{
		Path: param.Link,
		List: list,
	})
}
