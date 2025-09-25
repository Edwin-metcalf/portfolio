<script lang="ts">
    import { Lock, UserPlus, Vote, BookOpenText } from 'lucide-svelte';
    const projects = [
        {
            id: 1,
            title: "Action-map Political Information Website",
            description: "Worked in an agile team environment with multiple iterations. Used Ruby on Rails with MVC architecutre to build a wesbsite to learn about your representatives. Incorporated multiple API's to pull information along with omni auth2 for SSO. Also cached results onto a SQL database.",
            icon: Vote,
            tools: ["Ruby", "Rails", "SQL"]
        },
        {
            id: 2,
            title: "Secure File Sharing Software",
            description: "Built an end to end encypted secure file sharing Software in Go. Leveraged AES-CTR, HashKDF, HMACs, RSA signatures, and UUID's to ensure Integrity, Confidentiality and Authenticity.",
            icon: Lock,
            tools: ["Go"]
        },
        {
            id: 3,
            title: "This Website!",
            description: "Used Svelte, html, css, and Type Script. It is all hosted using Heroku",
            icon: UserPlus,
            tools: ["Svelte", "HTML", "CSS", "TS","Heroku"]
        },
        {
            id: 4,
            title: "Class Registration Software",
            description: "Created software that processes CSV files of students and their class selections along with classes and their openings. Then puts students in classes depending on factors like grade or draw number.",
            icon: BookOpenText,
            tools: ["Java"]
        }

    ]

</script>

<div class="portfolio-section" id="portfolio">
    <div class="portfolio-container">
        <h1 class="portfolio-title">Portfolio</h1>
        
        <div class="projects-grid">
            {#each projects as project}
                <div class="project-card">
                    <div class="card-inner">
                        <!-- Front of card -->
                        <div class="card-front">
                            <div class="card-icon">
                                <svelte:component this={project.icon} size={32} />
                            </div>
                            <h3 class="project-title">{project.title}</h3>
                        </div>
                        
                        <!-- Back of card -->
                        <div class="card-back">
                            <p class="project-description">{project.description}</p>
                            <div class="tech-stack">
                                {#each project.tools as tool}
                                    <span class="tech-button">{tool}</span>
                                {/each}
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</div>

<style>
    .portfolio-section {
        padding: 80px 5%;
        display: flex;
        justify-content: center;
    }
    
    .portfolio-container {
        max-width: 1200px;
        width: 100%;
    }
    
    .portfolio-title {
        font-size: 3rem;
        font-weight: 600;
        color: #f0f6fc;
        text-align: center;
        margin-bottom: 50px;
        position: relative;
    }
    
    .portfolio-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 3px;
        background: #00d4aa;
        margin: 20px auto;
        border-radius: 2px;
    }
    
    .projects-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 25px;
        justify-items: center;
    }
    
    .project-card {
        position: relative;
        width: 100%;
        max-width: 250px;
        min-height: 200px;
        cursor: pointer;
        perspective: 1000px;
    }
    
    .project-card:hover {
        transform: translateY(-5px);
    }
    
    .card-inner {
        position: relative;
        width: 100%;
        height: 100%;
        text-align: center;
        transition: transform 0.6s;
        transform-style: preserve-3d;
    }
    
    .project-card:hover .card-inner {
        transform: rotateY(180deg);
    }
    
    .card-front,
    .card-back {
        position: absolute;
        width: 100%;
        height: 200px;
        backface-visibility: hidden;
        background: rgba(33, 38, 45, 0.4);
        backdrop-filter: blur(10px);
        -webkit-backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.08);
        border-radius: 16px;
        padding: 30px 20px;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
    
    .card-back {
        transform: rotateY(180deg);
        background: rgba(13, 17, 23, 0.95);
        border-color: rgba(0, 212, 170, 0.3);
    }
    
    .card-icon {
        color: #00d4aa;
        margin-bottom: 20px;
        transition: transform 0.3s ease;
    }
    
    .project-title {
        font-size: 1.3rem;
        font-weight: 600;
        color: #f0f6fc;
        margin: 0;
        line-height: 1.3;
    }
    
    .project-description {
        color: #c9d1d9;
        font-size: 0.9rem;
        line-height: 1.5;
        margin-bottom: 20px;
        text-align: center;
    }
    
    .tech-stack {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        justify-content: center;
    }
    
    .tech-button {
        background: rgba(0, 212, 170, 0.15);
        color: #00d4aa;
        padding: 6px 12px;
        border-radius: 20px;
        font-size: 0.75rem;
        font-weight: 500;
        border: 1px solid rgba(0, 212, 170, 0.3);
        transition: all 0.2s ease;
    }
    
    .tech-button:hover {
        background: rgba(0, 212, 170, 0.25);
        transform: scale(1.05);
    }
</style>