<script lang="ts">
    import Courses from "./Courses.svelte";
    import { fly } from 'svelte/transition';
    import { quintOut } from 'svelte/easing';
    import{ onMount } from 'svelte';


    let visable: boolean = false;
    let sectionRef: HTMLElement | null = null;

    onMount(() => {
        if (!sectionRef) return;

        const handleIntersection: IntersectionObserverCallback = (entries, observer) => {
            for (const entry of entries){
                if (entry.isIntersecting) {
                    visable = true;
                    observer.unobserve(entry.target);
                }
            }
        };

        const options: IntersectionObserverInit = {
            threshold: 0.8,
        };

        const observer = new IntersectionObserver(handleIntersection, options);
        observer.observe(sectionRef);

        return () => observer.disconnect();
    });

</script>
<div class="about-section" id="about-me" bind:this={sectionRef}>
    <div class="about-container">
        {#if visable}
            <h1 class="about-title" in:fly={{ y: -50, duration: 1000, delay: 100, easing: quintOut}}>About Me</h1>
            <p class="about-text">
                Hello Welcome to my Website my name is Edwin, I am an aspiring full stack software developer. I have a strong love of problem solving and video games which continually sends me down rabbit holes of learning new frameworks or languages. Not only do I enjoy programming, I love history and am a double major at Vassar College. Through my curiosity I have spent time abroad studying in London and traveling Europe along with taking classes at UC Berkeley to not only expand my education but also for new experiences.
            </p>
            <Courses></Courses>
        {/if}
    </div>
</div>

<style>
    .about-section {
        padding: 80px 5%;
        display: flex;
        justify-content: center;
    }
    .about-container {
        text-align: center;
        max-width: 800px;
    }
    .about-title {
        font-size: 3rem;
        font-weight: 600;
        color: #f0f6fc;
        margin-bottom: 30px;
        position: relative;
    }
    .about-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 3px;
        background: #00d4aa;
        margin: 20px auto;
        border-radius: 2px;
    }
    .about-text {
        font-size: 1.2rem;
        line-height: 1.8;
        color: #c9d1d9;
        margin: 0;
    }

    @media (max-width: 768px) {
        .about-section {
            padding: 50px 5%;
        }
        
        .about-container {
            max-width: 100%;
        }
        
        .about-title {
            font-size: 2rem;
            margin-bottom: 20px;
        }
        
        .about-title::after {
            width: 50px;
            height: 2px;
            margin: 15px auto;
        }
        
        .about-text {
            font-size: 1rem;
            line-height: 1.6;
            text-align: left;
        }
    }
    
</style>