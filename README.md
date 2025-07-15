<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/tomatoCoderq/repeatro">
    <img src="https://avatars.githubusercontent.com/u/215998499?s=48&v=4" alt="Logo" width="80" height="80">
  </a>
  <h3 align="center">Repeatro – Anki-Style Vocabulary Learning App</h3>
  <p align="center">
    A modern web-based vocabulary learning tool inspired by Anki, built with Go and PostgreSQL.<br />
    <a href="https://github.com/tomatoCoderq/repeatro"><strong>Explore the docs »</strong></a>
    <br />
    <a href="https://github.com/tomatoCoderq/repeatro/issues">Issues</a>
    &middot;
    <a href="https://github.com/tomatoCoderq/repeatro/pulls">Pull Requests</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#built-with">Built With</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
    <li><a href="#contributors">Contributors</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

Repeatro is a modern, open-source vocabulary learning app inspired by Anki. It leverages spaced repetition (SM2 algorithm) to help users efficiently retain vocabulary. Organize your words into decks, track your progress, and enjoy a simple, effective learning experience.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Features

- Spaced repetition for efficient vocabulary retention (SM2 algorithm)
- JWT-based user authentication
- Decks to organize vocabulary by topic or language
- [In progress...] Language detection using [lingua-go][lingua-go]
- [In progress...] RESTful API with [Swaggo][swaggo] auto-generated Swagger docs
- PostgreSQL backend with [Goose][goose] for migrations

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Built With

- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Goose][goose] (migrations)
- [Swaggo][swaggo] (API docs)
- [lingua-go][lingua-go] (language detection)
- [JWT](https://jwt.io/) (authentication)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Getting Started

To get a local copy up and running, follow these steps.

### Prerequisites

- Go 1.18+
- PostgreSQL 15+
- [Goose][goose]

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/tomatoCoderq/repeatro.git
   cd repeatro
   ```
2. Install dependencies
   ```sh
   go mod tidy
   ```
3. Create the database
   ```sql
   CREATE DATABASE repeatro;
   ```
4. Set up Goose (see [Goose docs][goose] for details)
5. Configure your environment (see `config.example.toml` in the root)
6. Start the backend
   ```sh
   go run backend/card/cmd/card/main.go
   # or for other services, adjust the path accordingly
   ```
7. (Optional) Use [air](https://github.com/cosmtrek/air) for auto server restart

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

Once running, you can interact with the API (REST/gRPC) for deck and card management, user authentication, and spaced repetition review. See the [Swagger docs](backend/repeatro/docs/swagger.yaml) for API details.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Roadmap

- [ ] Import/export via CSV, JSON
- [ ] Enhance current stats
- [ ] Language detection
- [ ] RESTful API with Swagger docs
- [ ] Simple frontend

See the [open issues][issues-url] for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contact

Maintainer: [@tomatoCoderq](https://github.com/tomatoCoderq)

Project Link: [https://github.com/tomatoCoderq/repeatro](https://github.com/tomatoCoderq/repeatro)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Acknowledgments

- [Anki](https://apps.ankiweb.net/) – inspiration
- [Goose][goose]
- [Swaggo][swaggo]
- [lingua-go][lingua-go]
- [Img Shields](https://shields.io)
- [Best-README-Template](https://github.com/othneildrew/Best-README-Template)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributors

<a href="https://github.com/tomatoCoderq/repeatro/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=tomatoCoderq/repeatro" alt="Contributors" />
</a>

**Project contributors:**
- [@tomatoCoderq](https://github.com/tomatoCoderq)
- [@constabIe](https://github.com/constabIe)
- [@Kaghorz](https://github.com/Kaghorz)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/tomatoCoderq/repeatro.svg?style=for-the-badge
[contributors-url]: https://github.com/tomatoCoderq/repeatro/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/tomatoCoderq/repeatro.svg?style=for-the-badge
[forks-url]: https://github.com/tomatoCoderq/repeatro/network/members
[stars-shield]: https://img.shields.io/github/stars/tomatoCoderq/repeatro.svg?style=for-the-badge
[stars-url]: https://github.com/tomatoCoderq/repeatro/stargazers
[issues-shield]: https://img.shields.io/github/issues/tomatoCoderq/repeatro.svg?style=for-the-badge
[issues-url]: https://github.com/tomatoCoderq/repeatro/issues
[license-shield]: https://img.shields.io/github/license/tomatoCoderq/repeatro.svg?style=for-the-badge
[license-url]: https://github.com/tomatoCoderq/repeatro/blob/main/LICENSE
[goose]: https://github.com/pressly/goose
[swaggo]: https://github.com/swaggo/swag
[lingua-go]: https://github.com/pemistahl/lingua-go
