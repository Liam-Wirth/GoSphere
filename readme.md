# Simple implementation of a spinning cube in go, might expand it to more geometry in the future but for now it's just a lil cube
[![wakatime](https://wakatime.com/badge/user/d40f8d42-5a14-4981-a36e-39f7bd209ef3/project/b55fc834-68e0-4f29-a02a-0c693229d315.svg)](https://wakatime.com/badge/user/d40f8d42-5a14-4981-a36e-39f7bd209ef3/project/b55fc834-68e0-4f29-a02a-0c693229d315)
- First Time using Golang, was kinda fun, felt a bit different
- used other repos for reference to try and get an idea of how this works, this is in no way an original idea learned
    - namely this one:
    https://github.com/saatvikrao/Spinning-Cube/


# TODO
    [ ] LOD implementation
    [ ] look into offloading the computation for rendering to the GPU
    [ ] STL Parsing
    [ ] For the sphere/logic that will be extended for other shapes implement aspect ratio calculation based on given terminal dimensions right now I just guesstimate
    [ ] Refactor so that there is a unified "scene" element that items are placed within, this scene will manage lighting and stuff like that rn I'm just working on the sphere this is for the refactor
