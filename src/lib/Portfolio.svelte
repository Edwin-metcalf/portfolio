<script lang="ts">
    import { Lock, UserPlus, Vote, BookOpenText, X, Code, Database, Server, Blocks } from 'lucide-svelte';

    interface Project {
        id: number;
        title: string;
        description: string;
        icon: any;
        tools: string[];
        longDescription: string;
        highlights: string[];
    }
    
    const projects: Project[] = [
        {
            id: 1,
            title: "Action-map Political Information Website",
            description: "Action-map Political Inforamtion Website",
            icon: Vote,
            tools: ["Ruby", "Rails", "SQL"],
            longDescription: "Worked in an agile team environment with multiple iterations. Used Ruby on Rails with MVC architecutre to build a wesbsite to learn about your representatives. Incorporated multiple API's to pull information along with omni auth2 for SSO. Also cached results onto a SQL database.",
            highlights: ["Agile team environment to create a SaaS product ","Connected API calls and SQL Database", "built a responsive and intersting front end", "Used Cucumber and Rspec for exhaustive testing"]
        },
        {
            id: 2,
            title: "Secure File Sharing Software",
            description: "End to end encypted file sharing software using: ",
            icon: Lock,
            tools: ["GoLang"],
            longDescription: "Built an end to end encypted secure file sharing Software in Go. Leveraged AES-CTR, HashKDF, HMACs, RSA signatures, and UUID's to ensure Integrity, Confidentiality and Authenticity. Users are able to share and revoke access to files ",
            highlights: ["Employed Cyptographic schemes to securely store, share and send files","used both public key encyption and symetric key encyption"]
        },
        {
            id: 3,
            title: "This Website!",
            description: "My portfolio website",
            icon: UserPlus,
            tools: ["Svelte", "HTML", "CSS", "TS","Heroku"],
            longDescription: "Used Svelte, HTML, CSS, and Type Script. It is all hosted using Heroku",
            highlights: ["created an informative and good looking website "]
        },
        {
            id: 4,
            title: "Class Registration Software",
            description: "Class Registration Software",
            icon: BookOpenText,
            tools: ["Java"],
            longDescription: "Created software that processes CSV files of students and their class selections along with classes and their openings. Then puts students in classes depending on factors like grade or draw number.",
            highlights: ["Ability to load courses and students from CSV", "Return a CSV with all the students schedules with the correct classes enrolled"]
        }

    ];

    let selectedProject: Project|null = null;

    function openModal(project: Project): void {
        selectedProject = project;
        document.body.style.overflow = 'hidden';
    }
    function closeModal(): void {
        selectedProject = null;
        document.body.style.overflow = 'auto';
    }

    function handleKeydown(event: KeyboardEvent): void {
        if (event.key === 'Escape') {
            closeModal();
        }
    }

    function handleCardKeydown(event: KeyboardEvent, project: Project): void {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            openModal(project);
        }
    }

    function handleOverlayKeydown(event: KeyboardEvent): void {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            closeModal();
        }
    }

</script>

<svelte:window on:keydown={handleKeydown} />

<div class="portfolio-section" id="portfolio">
    <div class="portfolio-container">
        <h1 class="portfolio-title">Portfolio</h1>
        
        <div class="projects-grid">
            {#each projects as project}
                <div class="project-card" on:click={() => openModal(project)} on:keydown={(e) => handleCardKeydown(e, project)} role="button" tabindex="0" aria-label={`Open details for ${project.title}`}>
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
                                    {#if tool == "Java" || tool == "Ruby" || tool == "GoLang" || tool == "TS" || tool == "HTML" || tool == "CSS"}
                                        <span class="tech-button"> <Code size={15}></Code> {tool}</span>
                                    {:else if tool == "SQL"}
                                       <span class="tech-button"> <Database size={15}></Database> {tool}</span>
                                    {:else if tool == "Heroku"}
                                        <span class="tech-button"> <Server size={15}></Server> {tool}</span>
                                    {:else if tool == "Svelte" || tool == "Rails"}
                                        <span class="tech-button"> <Blocks size={15}></Blocks> {tool}</span>
                                    {:else}
                                        <span class="tech-button">{tool}</span>
                                    {/if}

                                {/each}
                            </div>
                            <span class="click-hint">Click for details →</span>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</div>

<!-- pop up time -->
{#if selectedProject}
    <div class="modal-overlay" on:click={closeModal} on:keydown={handleOverlayKeydown} role="button" tabindex="0" aria-label="Close modal">
        <div class="modal-content" on:click|stopPropagation on:keydown= {() => {}} role="dialog" aria-modal="true" aria-labelledby="modal-title" tabindex="0">
            <button class="modal-close" on:click={closeModal}>
                <X size={24}/>
            </button>
            <div class="modal-header">
                <div class="modal-icon">
                    <svelte:component this={selectedProject.icon} size={40}/>
                </div>
                <h2 class="modal-title">{selectedProject.title}</h2>
            </div>

            <div class="modal-body">
                <h3 class="section-title">Description</h3>
                <p class="modal-description"> {selectedProject.longDescription}</p>

                <div class="modal-section">
                    <h3 class="section-title">Highlights</h3>
                    <ul class="highlights-list">
                        {#each selectedProject.highlights as highlight}
                            <li>{highlight}</li>
                        {/each}
                    </ul>
                </div>
                <div class="modal-section">
                    <h3 class="section-title">Tech Stack</h3>
                    <ul class="tech-stack">
                        {#each selectedProject.tools as tool}
                            {#if tool == "Java" || tool == "Ruby" || tool == "GoLang" || tool == "TS" || tool == "HTML" || tool == "CSS"}
                                <span class="tech-button"> <Code size={15}></Code> {tool}</span>
                            {:else if tool == "SQL"}
                                <span class="tech-button"> <Database size={15}></Database> {tool}</span>
                            {:else if tool == "Heroku"}
                                <span class="tech-button"> <Server size={15}></Server> {tool}</span>
                            {:else if tool == "Svelte" || tool == "Rails"}
                                <span class="tech-button"> <Blocks size={15}></Blocks> {tool}</span>
                            {:else}
                                <span class="tech-button">{tool}</span>
                            {/if}
                        {/each}
                    </ul>
                </div>
            </div>
        </div>
    </div>
{/if}

<style>
    .portfolio-section {
        padding: 80px 5%;
        padding-left: 85px;
        display: flex;
        justify-content: center;
        min-height: 80vh;
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
        overflow: hidden;
        border-radius: 16px;
    }
    .project-card::before {
        content: "";
        position: absolute;
        top: 0;
        left: -75%;
        width: 50%;
        height: 100%;
        background: linear-gradient(
            120deg,
            rgba(0, 212, 170, 0) 0%,       
            rgba(0, 212, 170, 0.12) 50%,    
            rgba(0, 212, 170, 0) 100%       
        );
        transform: skewX(-25deg);
        animation: shine 4s infinite;
        pointer-events: none;
    }
    @keyframes shine {
        0% {
            left: -75%;
        }
        100% {
            left: 125%;
        }
    }
    .project-card:hover::before {
        animation: none;
        opacity: 0;
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
        top: 0;
        left: 0;
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
        box-sizing: border-box;
    }
    
    .card-back {
        transform: rotateY(180deg);
        background: rgba(13, 17, 23, 0.95);
        border-color: rgba(0, 212, 170, 0.3);
        box-shadow: 
            0 0 20px rgba(0, 212, 170, 0.3),
            0 0 40px rgba(0, 212, 170, 0.1),
            inset 0 0 20px rgba(0, 212, 170, 0.05);
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
    .click-hint {
        color: #00d4aa;
        font-size: 0.75rem;
        font-weight: 500;
        margin-top: 5px;
        opacity: 0.8;
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
        display: inline-flex;
        align-items: center;
        gap: 6px;
    }
    
    .tech-button:hover {
        background: rgba(0, 212, 170, 0.25);
        transform: scale(1.05);
    }

    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.8);
        backdrop-filter: blur(5px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 20px;
        animation: fadeIn 0.5s ease;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    .modal-content {
        background: rgba(13, 17, 23, 0.98);
        border: 1px solid rgba(0, 212, 170, 0.3);
        border-radius: 20px;
        max-width: 600px;
        width: 100%;
        max-height: 90vh;
        overflow-y: auto;
        position: relative;
        box-shadow: 
            0 0 20px rgba(0, 212, 170, 0.4),
            0 0 40px rgba(0, 212, 170, 0.2),
            inset 0 0 30px rgba(0, 212, 170, 0.05);
        animation: slideUp 0.3s ease;
    }
    @keyframes slideUp {
        from {
            transform: translateY(50px);
            opacity: 0;
        }
        to {
            transform: translateY(0);
            opacity: 1;
        }
    }
    .modal-close {
        position: absolute;
        top: 20px;
        right: 20px;
        background: rgba(255, 255, 255, 0.1);
        border: 1px solid rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        color: #f0f6fc;
        transition: all 0.3s ease;
        z-index: 1;
    }
    
    .modal-close:hover {
        background: rgba(0, 212, 170, 0.2);
        border-color: #00d4aa;
        transform: rotate(90deg);

    }

    .modal-header {
        padding: 40px 40px 20px;
        text-align: center;
        border-bottom: 4px solid rgba(255, 255, 255, 0.1);
    }
    .modal-icon {
        color: #00d4aa;
        margin-bottom: 15px;
    }

    .modal-title {
        font-size: 1.8rem;
        font-weight: 600;
        color: #f0f6fc;
        margin: 0;
        line-height: 1.3;
    }

    .modal-body {
        padding: 30px 40px 40px;
    }
    .modal-description {
        color: #c9d1d9;
        font-size: 1rem;
        line-height: 1.7;
        margin-bottom: 30px;
    }

    .modal-section {
        margin-bottom: 30px;
    }

    .section-title {
        font-size: 1.2rem;
        font-weight: 600;
        color: #00d4aa;
        margin-bottom: 15px;
    }

    .highlights-list {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    .highlights-list li {
        color: #c9d1d9;
        font-size: 0.95rem;
        line-height: 1.6;
        padding: 10px 0 10px 25px;
        position: relative;
    }

    .highlights-list li::before {
        content: '→';
        position: absolute;
        left: 0;
        color: #00d4aa;
        font-weight: bold;
    }

    @media (max-width: 768px) {
        .projects-grid {
            grid-template-columns: 1fr;
            gap: 20px;
        }
        
        .project-card {
            max-width: 100%;
            min-height: 180px;
        }
        
        .portfolio-title {
            font-size: 2.5rem;
        }

        .modal-content {
            max-height: 95vh;
            border-radius: 15px;
        }

        .modal-header,
        .modal-body {
            padding: 30px 25px;
        }

        .modal-title {
            font-size: 1.5rem;
        }
    }
</style>