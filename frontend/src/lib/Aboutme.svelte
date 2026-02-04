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
                Hi, I'm Edwin. I'm an entry-level full-stack software developer with a strong focus on building secure, well-designed systems. I enjoy problem solving, video games and love learning new frameworks or languages by building projects.
            </p>
            <p class="about-text">
                I have experience developing RESTful APIs and informative, user-friendly front ends, with an emphasis on security and clean system design. I also have hands-on experience with AI models and machine learning through natural language processing.             
            </p>

            <p class="about-text">
                Not only do I enjoy programming, I love history and am a double major at Vassar College. Through my curiosity I have spent time abroad studying in London and traveling Europe along with taking classes at UC Berkeley to not only expand my education but also for new experiences. These experiences have strengthened my adaptability, communication skills, and ability with all types of people. 
            </p>
            <p class="about-text">
                In addition to my academic and technical work, I am a varsity lacrosse player which has taught me leadership skills and reinforced the importance of teamwork, discipline, and working toward shared goals.
            </p>
            <div style="display: none;">
                <Courses></Courses>
            </div>
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
        margin-bottom: 0.5rem;
        text-align: left;
        text-indent: 8%;
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