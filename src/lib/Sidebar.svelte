<script lang="ts">
    import { onMount } from 'svelte';
    import { User, Code, Briefcase, Mail } from 'lucide-svelte';

    let visible = false;
    let activeSection = 'about-me'


    const sections = [
        {id: 'about-me', icon: User, label: 'About'},
        {id: 'skills', icon: Code, label: 'Skills'},
        {id: 'portfolio', icon: Briefcase, label: 'Portfolio'},
        {id: 'contact', icon: Mail, label: 'Contact'}
    ];

    onMount(() =>{
        const handleScroll = () => {

            visible = window.scrollY > 200;
            const scrollPosition = window.scrollY + 100;

            for(const section of sections) {
                    const element = document.getElementById(section.id);
                    if (element) {
                        const { offsetTop, offsetHeight } = element;
                        if (scrollPosition >= offsetTop && scrollPosition < offsetTop + offsetHeight) {
                            activeSection = section.id;
                            break;
                        }  
                    }

            }
        };

        window.addEventListener('scroll',handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    });

    function scrollToSection(sectionId: string) {
        const element = document.getElementById(sectionId);
        if (element) {
            element.scrollIntoView({behavior: 'smooth'});
        }
    }
</script>

<div class="sidebar" class:visible>
    <nav class="sidebar-nav">
        {#each sections as section}
            <button class="nav-item" class:active={activeSection === section.id} on:click={() => scrollToSection(section.id)} aria-label={`Navigate to ${section.label}`}>
                <div class="icon-wrapper">
                    <svelte:component this={section.icon} size={20}/>
                </div>
                <span class="nav-label">{section.label}</span>
            </button>

        {/each}
    </nav>
</div>

<style>
    .sidebar {
        position: fixed;
        left: 20px;
        top: 50%;
        transform: translateY(-50%) translateX(-120px); 
        z-index: 9999;
        transition: transform 0.4s ease, opacity 0.4s ease;
        opacity: 0;                                       
        pointer-events: none;                             
    }
    .sidebar.visible {
        transform: translateY(-50%) translateX(0);
        opacity: 1;                                      
        pointer-events: all;                             
    }
    .sidebar-nav {
        display: flex;
        flex-direction: column;
        gap: 12px;
        background: rgba(13, 17, 23, 0.95);
        backdrop-filter: blur(10px);
        -webkit-backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 16px;
        padding: 16px 12px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    }
    
    .nav-item {
        position: relative;
        display: flex;
        align-items: center;
        gap: 12px;
        background: transparent;
        border: none;
        cursor: pointer;
        padding: 12px;
        border-radius: 10px;
        transition: all 0.3s ease;
        color: #8b949e;
        overflow: hidden;
        white-space: nowrap;
    }
    
    .icon-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        min-width: 20px;
        transition: transform 0.3s ease;
    }
    
    .nav-label {
        font-size: 0.9rem;
        font-weight: 500;
        opacity: 0;
        max-width: 0;
        transition: all 0.3s ease;
        color: inherit;
    }
    
    .nav-item:hover {
        background: rgba(0, 212, 170, 0.1);
        color: #00d4aa;
        padding-right: 16px;
    }
    
    .nav-item:hover .nav-label {
        opacity: 1;
        max-width: 100px;
        margin-left: 4px;
    }
    
    .nav-item:hover .icon-wrapper {
        transform: scale(1.1);
    }
    
    .nav-item.active {
        background: rgba(0, 212, 170, 0.15);
        color: #00d4aa;
        border-left: 3px solid #00d4aa;
    }
    
    .nav-item.active .icon-wrapper {
        transform: scale(1.05);
    }
    
    /* Hide on mobile */
    @media (max-width: 768px) {
        .sidebar {
            display: none;
        }
    }

</style>