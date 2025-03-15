## Intro

This is a Basic Demo of the Ebitengine for the *Go* language.

I intend this demo to be brief and simply demonstrate to myself that this engine runs on my machine; 
though I might want to take things further with future branching out.

the end goal is for a solid set of experimentation to have taken place;


# wasm

in the folder with main.go in it; run 
```
>go run github.com/hajimehoshi/wasmserve@latest
```

# build to bin folder

go build -o bin/ main.go

## additional information:

    this is being written with my future self in mind;
    there is a great deal here that is grossly inefficient, and probably a lot that is bad practice;

    March 12, 2025: implemented the "Animated Sprite" class; noticed some substantial drops in framerate however; little bug or something not sure if it's because I'm making calls to timer or what;
        >>draw to a background image rather than the screen
            -tried this it seems to work pretty well;-march 13th, 2025

        ---
        other things March 13, 2025
        I'm incrasingly convinced that I would be better off restarting with a new project that would take what I've learned here and apply it to a better scope;
        


## TODO:


    The current size of the build is around 15 MB in size; assets only 6kbs 
    THIS IS despite it being incredibly primitive with barely any graphics, barely any gameplay, no levels, no AI;
    This is far too large to be a good browser game; -> need to figure out how to optimize to make it smaller, neater;
    IDEAS: perhaps chaging the 64 bit variables to more manageable 32 bit ones would help with this?
    other ideas: trim some of the fat with unused classes;

    other other ideas: additional feed 


so march 15th,

    moved to a new branch; got a lot of stuff done;
    the interesting thing is that there's a weird barrier with accessing the array of buttons in a button panel; such that there might need to be a total reworking of buttons in order to accomodate it;

    this is a valuable learning experience of course; I really really enjoy most of it; 

    the buttons and panels might be solved with an array of pointers pointing at the values stored in the array;
    