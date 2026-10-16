<script lang="ts">
    import { Lock, UserPlus, Vote, BookOpenText, X, Code, Database, Server, Blocks, BrainCircuit, FileCode2, Bot, Network, MessageSquare, Wand } from 'lucide-svelte';

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
            id: 7,
            title:  "Automated Evaluation for Vision-Language Models",
            description: "Isaacsim evaluation platform",
            icon: Bot,
            tools: ["Python", "Isaacsim"],
            longDescription: "Built a headless IsaacSim evaluation platform in Python for a AI model for evaluating Vision-Language models on autonomous control of a Franka Panda robot arm",
            highlights: ["automated execution of vision-language  manipulation tasks", "Translated natural language into robot actions", "automated evaluation loop including task execution, success/failure logging, and performance metrics"]
        },
        {
            id: 2,
            title: "Secure File Sharing Software",
            description: "End to end encypted file sharing software supporting secure storage, sharing, and revocation of user access",
            icon: Lock,
            tools: ["Go"],
            longDescription: "Built an end to end encypted secure file sharing Software in Go. Leveraged AES-CTR, HashKDF, HMACs, RSA signatures, and UUID's to ensure Integrity, Confidentiality and Authenticity. Users are able to share and revoke access to files ",
            highlights: ["Employed Cyptographic schemes to securely store, share and send files","used both public key encyption and symetric key encyption", "designed access-sharing and revoctation protocols"]
        },
        {
            id: 10,
            title: "MetCode",
            description: "Agentic command line assistant written in Go",
            icon: Wand,
            tools: ["Ollama", "tailscale", "Go"],
            longDescription: "MetCode is an agentic command line assistant that runs locally with Ollama, use with tailscale to run it on a PC from anywhere! MetCode provides a secure and private environment for executing commands. It can be used to automate tasks, generate code, and more.",
            highlights: ["Agentic command line assistant", "Runs locally with Ollama", "Safe Secure and private"]
        },
        {
            id: 8,
            title: "Bit Torrent Client",
            description: "Concurrent BitTorrent-style client using POSIX UDP",
            icon: Network,
            tools: ["C", "UDP", "Socket Programming"],
            longDescription: "Implemented concurrent peer downloads and sliding-window reliability protocol on top of UDP for high throughput file transfer",
            highlights: ["POSIX UDP socket programming", "sliding window reliability on top of UDP protocol", "Download chunks of files concurrently from multiple peers with SHA-1 Verficaiton"]
        },
        {
            id: 3,
            title: "This Website!",
            description: "My portfolio website",
            icon: UserPlus,
            tools: ["Svelte", "HTML", "CSS", "TS","Heroku","Go", "SQL"],
            longDescription: "Designed and deployed a fill stack personal portfolio using Svelte/TypeScript and a Go/PostgeSQL backend. Global leaderboard for custom game. Clash Royale matchup generator and more.",
            highlights: ["created an informative and good looking website", "Custom space invaders/vampire survivors game with global leaderboard", "Dota 2 and Clash royale game analysis", "Clash royale matchup generator"]
        }, 
        {
            id: 1,
            title: "Action-map Political Information Website",
            description: "Action-map Political Inforamtion Website",
            icon: Vote,
            tools: ["Ruby", "Rails", "SQL"],
            longDescription: "Ruby on Rails web application with MVC architecutre to help users learn about elected representatives. Incorporated external API's to pull information along with omni auth2 for SSO and saved information for users to see on SQL database.",
            highlights: ["integrated Google Civic information and external APIs", "Devleoped MVC components and responsive and intersting front end", "Used Cucumber and Rspec for testing", "Agile development"]
        },
        {
            id: 6,
            title: "Compiler",
            description: "Compiled custom programming language to RISC-V",
            icon: FileCode2,
            tools: ["Ocaml", "RISC-V"],
            longDescription: "Designed and implemented a compiler and a parser that used descent parsing and AST construction to generate RISC-V assembly from a custom .331 programming language ",
            highlights: ["recursive descent parsing and abstract syntax tree construction", "type checking for expressions, functions, conditions, and control flow", "optimizations: loop-invariant hoisting, constant fold, strength reduction, copy propagation, etc"]
        },
        {
            id: 1,
            title: "Action-map Political Information Website",
            description: "Action-map Political Inforamtion Website",
            icon: Vote,
            tools: ["Ruby", "Rails", "SQL"],
            longDescription: "Ruby on Rails web application with MVC architecutre to help users learn about elected representatives. Incorporated external API's to pull information along with omni auth2 for SSO and saved information for users to see on SQL database.",
            highlights: ["integrated Google Civic information and external APIs", "Devleoped MVC components and responsive and intersting front end", "Used Cucumber and Rspec for testing", "Agile development"]
        },
        {
            id: 5,
            title: "NLP paper: Fan Sentiment Correlation with MLB team success",
            description: "Temporal Analysis of sentiment of MLB fans",
            icon: BrainCircuit,
            tools: ["Python", "Scikit","Pandas"],
            longDescription: "Conducted research on sentiment analysis. We used VADER classifier to train a logistic regression model that was then used to conduct a temporal analysis of the SF Giants and the NY Yankees fanbases.",
            highlights: ["finding, processing, and cleaning of reddit data", "logistic regression model trained and tested", "temporal analysis of correlation of teams success and fan sentiment","Wrote a paper showing our results of Yankee fans being more volitile through the 2017 season"]
        },
        {
            id: 9,
            title: "Internet Relay Chat client",
            description: "Client written with POSIX library in C leverages TCP",
            icon: MessageSquare,
            tools: ["C", "TCP", "Socket Programming"],
            longDescription: "Implemented a client that is in accordance with the IRC protocol with mutiple commnds",
            highlights: ["POSIX TCP socket programming", "join, part, quit, server, msg, nick, whois and more commands"]
        },
        {
            id: 4,
            title: "Class Registration Software",
            description: "Class Registration Software",
            icon: BookOpenText,
            tools: ["Java"],
            longDescription: "Created software that processes CSV files of students and their class selections along with classes and their openings. Then puts students in classes depending on factors like grade or draw number.",
            highlights: ["Ability to load courses and students from CSV", "Return a CSV with all the students schedules with the correct classes enrolled"]
        },
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
                                    {#if tool == "Java" || tool == "Ruby" || tool == "GoLang" || tool == "TS" || tool == "HTML" || tool == "CSS" || tool == "Go"}
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
        color: var(--text);
        text-align: center;
        margin-bottom: 50px;
        position: relative;
        font-family: var(--font-display);
    }
    
    .portfolio-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 1px;
        background: var(--border);
        margin: 20px auto;
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
        max-width: 275px;
        min-height: 200px;
        cursor: pointer;
        border-radius: var(--radius);
    }
    .card-inner {
        position: relative;
        width: 100%;
        height: 100%;
        text-align: center;
    }
    
    .card-front,
    .card-back {
        width: 100%;
        height: 200px;
        background: var(--panel);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 30px 20px;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        box-sizing: border-box;
        transition: transform 0.15s ease, border-color 0.15s ease;
    }
    
    .card-back {
        display: none;
    }
    
    .project-card:hover .card-front {
        transform: translateY(-2px);
        border-color: var(--mint);
    }
    
    .card-icon {
        color: var(--mint);
        margin-bottom: 20px;
    }
    
    .project-title {
        font-size: 1.2rem;
        font-weight: 600;
        color: var(--text);
        margin: 0;
        line-height: 1.3;
        font-family: var(--font-display);
    }
    
    .project-description {
        color: var(--text-muted);
        font-size: 0.9rem;
        line-height: 1.5;
        margin-bottom: 20px;
        text-align: center;
    }
    .click-hint {
        color: var(--mint);
        font-size: 0.75rem;
        font-weight: 500;
        margin-top: 5px;
    }
    
    .tech-stack {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        justify-content: center;
    }
    
    .tech-button {
        background: transparent;
        color: var(--mint);
        padding: 6px 12px;
        border-radius: var(--radius);
        font-size: 0.75rem;
        font-weight: 500;
        border: 1px solid var(--border);
        transition: all 0.15s ease;
        display: inline-flex;
        align-items: center;
        gap: 6px;
        font-family: var(--font-mono);
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
        background: var(--panel);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        max-width: 600px;
        width: 100%;
        max-height: 90vh;
        overflow-y: auto;
        position: relative;
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
        background: transparent;
        border: 1px solid var(--border);
        border-radius: var(--radius);
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        color: var(--text);
        transition: all 0.15s ease;
        z-index: 1;
    }
    
    .modal-close:hover {
        background: rgba(62, 207, 192, 0.15);
        border-color: var(--mint);
    }

    .modal-header {
        padding: 40px 40px 20px;
        text-align: center;
        border-bottom: 1px solid var(--border);
    }
    .modal-icon {
        color: var(--mint);
        margin-bottom: 15px;
    }

    .modal-title {
        font-size: 1.8rem;
        font-weight: 600;
        color: var(--text);
        margin: 0;
        line-height: 1.3;
        font-family: var(--font-display);
    }

    .modal-body {
        padding: 30px 40px 40px;
    }
    .modal-description {
        color: var(--text-muted);
        font-size: 1rem;
        line-height: 1.7;
        margin-bottom: 30px;
    }

    .modal-section {
        margin-bottom: 30px;
    }

    .section-title {
        font-size: 1rem;
        font-weight: 600;
        color: var(--mint);
        margin-bottom: 15px;
        font-family: var(--font-mono);
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .highlights-list {
        list-style: none;
        padding: 0;
        margin: 0;
    }

    .highlights-list li {
        color: var(--text);
        font-size: 0.95rem;
        line-height: 1.6;
        padding: 10px 0;
        border-bottom: 1px solid var(--border);
    }

    .highlights-list li:last-child {
        border-bottom: none;
    }

    @media (max-width: 768px) {
        .portfolio-section {
            padding: 80px 5%;

        }
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