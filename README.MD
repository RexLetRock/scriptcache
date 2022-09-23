[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="http://github.com">
    <img src="extra/images/logo.png" alt="Logo" width="80" height="80"
  </a>

  <h3 align="center">SCRIPTCACHE</h3>

  <p align="center">
    SCRIPTCACHE - CACHE ALLOW SCRIPTING TO WRITE/LOAD VIA DATABASE
    <br />
    <a href="https://github.com/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/">View Demo</a>
    ·
    <a href="https://github.com/">Report Bug</a>
    ·
    <a href="https://github.com/">Request Feature</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->
<details open="open">
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
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

<!-- [![Product Name Screen Shot][product-screenshot]](https://example.com) -->
ZLIB help go code with safe and simple !!!

Here's why:
* Your time should be focused on creating something amazing. A project that solves a problem and helps others
* You shouldn't be doing the same tasks over and over like creating a SPA from scratch
* You should element DRY principles to the rest of your life :smile:

A list of commonly used resources that I find helpful are listed in the acknowledgements.

### Built With

This section should list bItem major frameworks that you built your project using. Leave bItem add-ons/plugins for the acknowledgements section. Here are a few examples.
* [Golang](https://go.dev)


<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* golang

### Installation
sysctl -w fs.file-max=2000500
sysctl -w fs.nr_open=2000500
sysctl -w net.nf_conntrack_max=2000500
ulimit -n 2000500

sysctl -w net.ipv4.tcp_tw_recycle=1
sysctl -w net.ipv4.tcp_tw_reuse=1

1. Get golang, go-swagger :
  ```sh
  go get -u github.com/swaggo/swag/cmd/swag
  wget https://github.com/go-swagger/go-swagger/releases/download/v0.28.0/swagger_linux_amd64 > swag
  ```
2. Clone repo
  ```sh
  git clone https://github.com
  ```
3. Init swagger document
  ```sh
  make swag
  ```
4. Run server
  ```sh
  make run
  ```

5. Open url : http://localhost:3000/swagger/index.html

<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_


<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/) for a list of proposed features (and known issues).


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. bItem contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/ZFeature`)
3. Commit your Changes (`git commit -m 'Add some ZFeature'`)
4. Push to the Branch (`git push origin feature/ZFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under Copyright License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Your Name - [@RexLetRock](https://twitter.com/RexLetRock) - RexLetRock@gmail.com
Project Link - [https://github.com/RexLetRock/zlib](https://github.com/RexLetRock/zlib)



<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements
* [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
* [Img Shields](https://shields.io)
* [Choose an Open Source License](https://choosealicense.com)
* [GitHub Pages](https://pages.github.com)
* [Animate.css](https://daneden.github.io/animate.css)
* [Loaders.css](https://connoratherton.com/loaders)
* [Slick Carousel](https://kenwheeler.github.io/slick)
* [Smooth Scroll](https://github.com/cferdinandi/smooth-scroll)
* [Sticky Kit](http://leafo.net/sticky-kit)
* [JVectorMap](http://jvectormap.com)
* [Font Awesome](https://fontawesome.com)





<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
<!-- ?style=for-the-badge -->
[contributors-shield]: https://img.shields.io/github/contributors/RexLetRock/zlib.svg
[contributors-url]: https://github.com/RexLetRock/zlib/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/RexLetRock/zlib.svg
[forks-url]: https://github.com/RexLetRock/zlib/network/members
[stars-shield]: https://img.shields.io/github/stars/RexLetRock/zlib.svg
[stars-url]: https://github.com/RexLetRock/zlib/stargazers
[issues-shield]: https://img.shields.io/github/issues/RexLetRock/zlib.svg
[issues-url]: https://github.com/RexLetRock/zlib/issues
[license-shield]: https://img.shields.io/github/license/RexLetRock/zlib.svg
[license-url]: https://github.com/RexLetRock/zlib/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/rexletrock
[product-screenshot]: images/screenshot.png

ssh -i ~/keys/doprivatekey -L 59000:localhost:5901 -C -N -l root 14.225.17.201

KUBERNETS
ssh -i ~/keys/doprivatekey -L 59000:localhost:8001 -C -N -l root 14.225.17.201
http://localhost:59000/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/pod?namespace=_all

MASTER
eyJhbGciOiJSUzI1NiIsImtpZCI6IjhfNXhudzYwN1FkQkRyX1hRYjVwbnUtNk1tZUV0TEsyZHRHZUVYMm9aU00ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi11c2VyLXRva2VuLXJzOXA1Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImFkbWluLXVzZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI0MDIyMzdiMC1kYzY3LTQ3MDAtYjg1Ni1jY2FmODJkNTg5Y2IiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZXJuZXRlcy1kYXNoYm9hcmQ6YWRtaW4tdXNlciJ9.H4fD51UQoF9dAZJTINNFjnrerb8Ov5LbQYvQQyBCez4LDsDod7iUw-fv5N7To9_MVwfWbwuF8zUnw7T8nS8cL-5mAH9JJFe6Ic0pxwsvZ6iv5KkmZ_pzkkTNc3XVjSWyAuQCLX4IDZB8riSOiodwLrOPAq4nCt7bJgGlphH0V595zLY1UoHN59snEeFiCi5St7PAwj94RAkhiEKx8YJoJs8qO5MIJz5D2MlbV1muGVHBtXf5K3GjNVCGKe-XQewH0RKHGoFNB28WRL8OH-9GN2bXNc4rDml8TASPgL3puiSZj98H1c5PIWsjRwVLSw8GBsNYxuBf8QNdoG_WmdK_Jw

READ
eyJhbGciOiJSUzI1NiIsImtpZCI6IjhfNXhudzYwN1FkQkRyX1hRYjVwbnUtNk1tZUV0TEsyZHRHZUVYMm9aU00ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJyZWFkLW9ubHktdXNlci10b2tlbi1nYnR6aCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJyZWFkLW9ubHktdXNlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImFkMmE0MTZmLWQzYTEtNGQ1ZC1iNDdlLTdkZmNjY2YxOTJhMyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlcm5ldGVzLWRhc2hib2FyZDpyZWFkLW9ubHktdXNlciJ9.iBKZWFeMjXPLxsx4azvwfquZ8e_NIDyziciWz-PzKoQuMdIGiiqwoCcMx9ee9jpNZZvoKq2Q9F1To8p2x9nrMgCkdOaLFTm7HE_6iJULl_dWevv6c7xZNO1_1_LPwcd0fJ_mAS3Zy5plWWH4j8MrupMzdL9kQYun8fI7bPGBZJqxSYBLdE6xCf2-o28kAPCTNzUmqvoXyMaiceTOOT2I94mNapHeYAi39qIk_TM2VtV7Xv5WXVh4N6ZUVM92f_OMVdt7O-DsIrf2jTVCkGPK7YlS3pmpjxYyNY7U-nyKdDCMo0pqMOcUkdKwcOcWeLdivaNnPb8EoABW-BOEfabcXw

KUBE SERVER
14.225.17.201 (kube master)
14.225.17.202 -p3333 (kube services)
master -> ssh root@192.168.1.101 (kube nodes)

KUBE DASHBOARD
Tunnel : ssh -i ~/keys/doprivatekey -L 59000:localhost:8001 -C -N -l root 14.225.17.201
Weblink: http://localhost:59000/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/pod?namespace=_all