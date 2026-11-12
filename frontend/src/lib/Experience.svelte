<script lang="ts">
    import { onMount } from 'svelte';

    interface Role {
        period: string;
        title: string;
        org: string;
        summary?: string;
        bullets?: string[];
        tags?: string[];
    }

    // Placeholder dates and copy — replace with the real details.
    const history: Role[] = [
        {
            period: '2017 — 2021',
            title: 'Soccer Referee',
            org: 'Middle & high school',
            summary: 'Officiated middle and high school matches, managing game flow and player safety.'
        },
        {
            period: '2021 — 2022',
            title: 'Lacrosse Coach',
            org: 'Youth lacrosse',
            summary: 'Coached youth lacrosse, running practices and teaching fundamentals.'
        },
        {
            period: 'Summer 2022',
            title: 'Server',
            org: 'Independent wine & beer bar',
            summary: 'Guided guests through a rotating list of independent wines and beers.'
        },
        {
            period: 'Summer 2023',
            title: 'Market Research Intern',
            org: 'Rainmaker Capital',
            bullets: [
                'Built market and competitor research for the investment team.',
                'Turned findings into concise write-ups the team could act on.'
            ]
        },
        {
            period: '2024 — 2025',
            title: 'Software Engineering Intern',
            org: 'Rainmaker Capital',
            bullets: [
                'Built and shipped internal tooling used by the investment team.',
                'Worked across the stack on data-heavy features and APIs.',
                'Wrote tests and reviewed code alongside the engineering team.'
            ],
            tags: ['Python', 'TypeScript', 'SQL']
        }
    ];

    const present: Role[] = [
        {
            period: 'Present',
            title: 'Software Engineer',
            org: 'DataAnnotation (contract)',
            summary: 'Contract software engineer building and reviewing code used to train AI models.',
            bullets: [
                'Build and review code across languages and frameworks for model training.',
                'Evaluate model outputs and write detailed feedback for improvement.'
            ],
            tags: ['Python', 'TypeScript']
        },
        {
            period: 'Present',
            title: 'Founder & Maintainer',
            org: 'MetCode (open source)',
            summary: 'Open-source agentic command-line assistant written in Go.',
            bullets: [
                'Runs locally with Ollama, with Tailscale for access from anywhere.',
                'Focused on a private, secure environment for agentic tasks.'
            ],
            tags: ['Go', 'Ollama', 'Tailscale']
        }
    ];

    let headingVisible = false;
    let forkVisible = false;
    let sectionRef: HTMLElement | null = null;
    let forkRef: HTMLElement | null = null;

    onMount(() => {
        const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
        if (reduceMotion) {
            headingVisible = true;
            forkVisible = true;
            return;
        }

        const observer = new IntersectionObserver(
            (entries) => {
                for (const entry of entries) {
                    if (!entry.isIntersecting) continue;
                    if (entry.target === sectionRef) headingVisible = true;
                    if (entry.target === forkRef) forkVisible = true;
                    observer.unobserve(entry.target);
                }
            },
            { threshold: 0.15 }
        );

        if (sectionRef) observer.observe(sectionRef);
        if (forkRef) observer.observe(forkRef);

        return () => observer.disconnect();
    });
</script>

<div class="experience-section" id="experience" bind:this={sectionRef}>
    <div class="experience-container">
        <h1 class="experience-title" class:visible={headingVisible}>Experience</h1>

        <div class="timeline">
            {#each history as role}
                <article class="timeline-entry">
                    <span class="node" aria-hidden="true"></span>
                    <div class="entry-date">{role.period}</div>
                    <div class="entry-content">
                        <h3 class="entry-title">{role.title}</h3>
                        <p class="entry-org">{role.org}</p>
                        {#if role.summary}
                            <p class="entry-summary">{role.summary}</p>
                        {/if}
                        {#if role.bullets?.length}
                            <ul class="entry-bullets">
                                {#each role.bullets ?? [] as bullet}
                                    <li>{bullet}</li>
                                {/each}
                            </ul>
                        {/if}
                        {#if role.tags?.length}
                            <div class="entry-tags">
                                {#each role.tags ?? [] as tag}
                                    <span class="tag">{tag}</span>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </article>
            {/each}

            <div class="fork-row" bind:this={forkRef}>
                <div class="fork-date">Present</div>
                <div class="fork-canvas" class:visible={forkVisible}>
                    <svg class="fork-desktop" viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true">
                        <path class="branch" d="M0 0 C0 52 24.2 38 24.2 100" pathLength="1" />
                        <path class="branch" d="M0 0 C0 52 75.8 38 75.8 100" pathLength="1" />
                    </svg>
                    <svg class="fork-mobile" viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true">
                        <path class="branch" d="M0 0 L0 100" pathLength="1" />
                    </svg>
                    <span class="junction" aria-hidden="true"></span>
                </div>
            </div>

            <div class="present">
                {#each present as role}
                    <article class="current-role">
                        <h3 class="entry-title">{role.title}</h3>
                        <p class="entry-org">{role.org}</p>
                        {#if role.summary}
                            <p class="entry-summary">{role.summary}</p>
                        {/if}
                        {#if role.bullets?.length}
                            <ul class="entry-bullets">
                                {#each role.bullets ?? [] as bullet}
                                    <li>{bullet}</li>
                                {/each}
                            </ul>
                        {/if}
                        {#if role.tags?.length}
                            <div class="entry-tags">
                                {#each role.tags ?? [] as tag}
                                    <span class="tag">{tag}</span>
                                {/each}
                            </div>
                        {/if}
                    </article>
                {/each}
            </div>
        </div>
    </div>
</div>

<style>
    .experience-section {
        padding: 80px 5%;
        display: flex;
        justify-content: center;
    }

    .experience-container {
        max-width: 900px;
        width: 100%;
    }

    .experience-title {
        font-size: 3rem;
        font-weight: 600;
        color: var(--text);
        text-align: center;
        margin-bottom: 50px;
        position: relative;
        font-family: var(--font-display);
        opacity: 0;
        transform: translateY(-20px);
        transition: opacity 0.7s ease, transform 0.7s cubic-bezier(0.22, 1, 0.36, 1);
    }

    .experience-title.visible {
        opacity: 1;
        transform: none;
    }

    .experience-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 1px;
        background: var(--border);
        margin: 20px auto;
    }

    .timeline {
        width: 100%;
    }

    .timeline-entry {
        position: relative;
        display: grid;
        grid-template-columns: 150px 1fr;
    }

    .timeline-entry::before {
        content: '';
        position: absolute;
        left: 150px;
        top: 0;
        bottom: 0;
        width: 1px;
        background: var(--border);
    }

    .entry-date {
        grid-column: 1;
        text-align: right;
        padding: 3px 28px 0 0;
        font-family: var(--font-mono);
        font-size: 0.75rem;
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--text-muted);
        line-height: 1.5;
    }

    .node {
        position: absolute;
        left: 150px;
        top: 6px;
        width: 9px;
        height: 9px;
        border-radius: 50%;
        background: var(--bg);
        border: 1px solid var(--border);
        transform: translateX(-50%);
        transition: background 0.15s ease, border-color 0.15s ease;
        z-index: 1;
    }

    .timeline-entry:hover .node {
        background: var(--mint);
        border-color: var(--mint);
    }

    .entry-content {
        grid-column: 2;
        padding: 0 0 44px 30px;
    }

    .entry-title {
        margin: 0;
        font-family: var(--font-display);
        font-size: 1.35rem;
        font-weight: 600;
        color: var(--text);
        line-height: 1.3;
    }

    .entry-org {
        margin: 5px 0 0;
        font-family: var(--font-mono);
        font-size: 0.75rem;
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--mint);
    }

    .entry-summary {
        margin: 10px 0 0;
        color: var(--text-muted);
        font-size: 0.95rem;
        line-height: 1.65;
        max-width: 62ch;
    }

    .entry-bullets {
        list-style: none;
        margin: 12px 0 0;
        padding: 0;
        max-width: 62ch;
    }

    .entry-bullets li {
        position: relative;
        padding-left: 16px;
        margin-bottom: 6px;
        color: var(--text-muted);
        font-size: 0.95rem;
        line-height: 1.6;
    }

    .entry-bullets li:last-child {
        margin-bottom: 0;
    }

    .entry-bullets li::before {
        content: '';
        position: absolute;
        left: 0;
        top: 0.72em;
        width: 6px;
        height: 1px;
        background: var(--text-muted);
    }

    .entry-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        margin-top: 14px;
    }

    .tag {
        font-family: var(--font-mono);
        font-size: 0.7rem;
        color: var(--mint);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 4px 9px;
    }

    .fork-row {
        display: grid;
        grid-template-columns: 150px 1fr;
    }

    .fork-date {
        grid-column: 1;
        text-align: right;
        padding: 0 28px 0 0;
        font-family: var(--font-mono);
        font-size: 0.75rem;
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--mint);
        line-height: 1;
        transform: translateY(-50%);
    }

    .fork-canvas {
        grid-column: 2;
        position: relative;
        height: 104px;
    }

    .fork-canvas svg {
        display: block;
        width: 100%;
        height: 100%;
        overflow: visible;
    }

    .fork-mobile {
        display: none;
    }

    .branch {
        fill: none;
        stroke: var(--mint);
        stroke-width: 1;
        stroke-opacity: 0.55;
        vector-effect: non-scaling-stroke;
        stroke-dasharray: 1;
        stroke-dashoffset: 1;
        transition: stroke-dashoffset 1s cubic-bezier(0.22, 1, 0.36, 1) 0.1s;
    }

    .fork-canvas.visible .branch {
        stroke-dashoffset: 0;
    }

    .junction {
        position: absolute;
        left: 0;
        top: 0;
        width: 9px;
        height: 9px;
        border-radius: 50%;
        background: var(--mint);
        transform: translate(-50%, -50%);
        box-shadow: 0 0 0 4px var(--bg);
    }

    .present {
        margin-left: 150px;
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 24px;
    }

    .current-role {
        position: relative;
        background: var(--panel);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 26px 24px 24px;
        transition: border-color 0.15s ease, transform 0.15s ease;
    }

    .current-role::before {
        content: '';
        position: absolute;
        top: -5px;
        left: 50%;
        width: 9px;
        height: 9px;
        border-radius: 50%;
        background: var(--mint);
        transform: translateX(-50%);
        box-shadow: 0 0 0 4px var(--bg);
    }

    .current-role:hover {
        border-color: var(--mint);
        transform: translateY(-2px);
    }

    .current-role .entry-title {
        font-size: 1.2rem;
    }

    @media (max-width: 768px) {
        .experience-section {
            padding: 50px 5%;
        }

        .experience-title {
            font-size: 2rem;
            margin-bottom: 30px;
        }

        .experience-title::after {
            width: 50px;
        }

        .timeline-entry {
            display: block;
            padding-left: 22px;
        }

        .timeline-entry::before {
            left: 0;
        }

        .node {
            left: 0;
            top: 5px;
        }

        .entry-date {
            text-align: left;
            padding: 0 0 8px;
        }

        .entry-content {
            padding: 0 0 36px;
        }

        .fork-row {
            display: block;
        }

        .fork-date {
            text-align: left;
            padding: 0 0 8px 22px;
            transform: none;
        }

        .fork-canvas .fork-desktop {
            display: none;
        }

        .fork-canvas .fork-mobile {
            display: block;
        }

        .fork-canvas {
            height: 36px;
        }

        .present {
            position: relative;
            margin-left: 0;
            padding-left: 22px;
            grid-template-columns: 1fr;
            gap: 20px;
        }

        .present::before {
            content: '';
            position: absolute;
            left: 0;
            top: 0;
            bottom: 0;
            width: 1px;
            background: linear-gradient(to bottom, var(--border) 0%, var(--border) 82%, transparent 100%);
        }

        .current-role::before {
            top: 26px;
            left: -22px;
            transform: translate(-50%, -50%);
        }
    }
</style>
