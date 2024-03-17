<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]
![API Coverage][APICoverage-shield]




<!-- PROJECT LOGO -->
<br />
<div align="center">
  <!-- <a href="https://github.com/azaurus1/go-pinot-api">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a> -->

<h3 align="center">go-pinot-api</h3>

  <p align="center">
    A Go Library for interacting with the Apache Pinot Controller
    <br />
    <a href="https://github.com/azaurus1/go-pinot-api"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <!-- <a href="https://github.com/azaurus1/go-pinot-api">View Demo</a>
    · -->
    <a href="https://github.com/azaurus1/go-pinot-api/issues">Report Bug</a>
    ·
    <a href="https://github.com/azaurus1/go-pinot-api/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li>
    <a href="#usage">Usage</a>
    <ul>
        <li><a href="#creating-a-user">Creating a user</a></li>
        <li><a href="#creating-a-schema">Creating a schema</a></li>
    </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

A library for interacting with Apache Pinot Controllers via the Controller REST API.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Go][Go]][Go-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started
### Installation

```sh
go get github.com/azaurus1/go-pinot-api
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

### Creating a user:
```go
user := pinotModel.User{
  Username:  "user",
  Password:  "password",
  Component: "BROKER",
  Role:      "admin",
}

userBytes, err := json.Marshal(user)
if err != nil {
  log.Panic(err)
}

// Create User
createResp, err := client.CreateUser(userBytes)
if err != nil {
  log.Panic(err)
}
```

### Creating a schema:
```go
f, err := os.Open(schemaFilePath)
if err != nil {
  log.Panic(err)
}

defer f.Close()

var schema pinotModel.Schema
err = json.NewDecoder(f).Decode(schema)
if err != nil {
  log.Panic(err)
}

_, err = client.CreateSchema(schema)
if err != nil {
	log.Panic(err)
}

```

_For more examples, please refer to the [Documentation](https://example.com)_



<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] User Management
- [ ] Schema Management
- [ ] Table Management
- [ ] Segment Management
- [ ] Tenant Management
- [ ] Cluster Management
- [ ] Instance Management
- [ ] Task Management

See the [open issues](https://github.com/azaurus1/go-pinot-api/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

See the [Contributing Guide]()

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Liam Aikin - liamaikin@gmail.com

Project Link: [https://github.com/azaurus1/go-pinot-api](https://github.com/azaurus1/go-pinot-api)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* Many thanks to [Hagen H (Hendoxc)](https://github.com/hendoxc)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/azaurus1/go-pinot-api.svg?style=for-the-badge
[contributors-url]: https://github.com/azaurus1/go-pinot-api/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/azaurus1/go-pinot-api.svg?style=for-the-badge
[forks-url]: https://github.com/azaurus1/go-pinot-api/network/members
[stars-shield]: https://img.shields.io/github/stars/azaurus1/go-pinot-api.svg?style=for-the-badge
[stars-url]: https://github.com/azaurus1/go-pinot-api/stargazers
[issues-shield]: https://img.shields.io/github/issues/azaurus1/go-pinot-api.svg?style=for-the-badge
[issues-url]: https://github.com/azaurus1/go-pinot-api/issues
[license-shield]: https://img.shields.io/github/license/azaurus1/go-pinot-api.svg?style=for-the-badge
[license-url]: https://github.com/azaurus1/go-pinot-api/blob/main/LICENSE
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/liam-aikin
[product-screenshot]: images/screenshot.png
[Next.js]: https://img.shields.io/badge/next.js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white
[Next-url]: https://nextjs.org/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Vue.js]: https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vuedotjs&logoColor=4FC08D
[Vue-url]: https://vuejs.org/
[Angular.io]: https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white
[Angular-url]: https://angular.io/
[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00
[Svelte-url]: https://svelte.dev/
[Laravel.com]: https://img.shields.io/badge/Laravel-FF2D20?style=for-the-badge&logo=laravel&logoColor=white
[Laravel-url]: https://laravel.com
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com 
[Go-url]: https://go.dev
[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[APICoverage-shield]: https://img.shields.io/badge/dynamic/json?url=https://api-coverage-server.azaurus.dev/badge&query=coverage&label=API%20Coverage&color=blue&suffix=%25&style=for-the-badge
