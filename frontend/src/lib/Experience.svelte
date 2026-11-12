<script lang="ts">
    import { onMount } from 'svelte';

    type Weight = 'recent' | 'mid' | 'early';

    interface Role {
        period: string;
        title: string;
        org: string;
        weight?: Weight;
        summary?: string;
        bullets?: string[];
        tags?: string[];
    }

    // Placeholder dates and copy — replace with the real details.
    const history: Role[] = [
        {
            period: '2024 — 2025',
            title: 'Software Engineering Intern',
            org: 'Rainmaker Capital',
            weight: 'recent',
            bullets: [
                'Built and shipped internal tooling used by the investment team.',
                'Worked across the stack on data-heavy features and APIs.',
                'Wrote tests and reviewed code alongside the engineering team.'
            ],
            tags: ['Python', 'TypeScript', 'SQL']
        },
        {
            period: 'Summer 2023',
            title: 'Market Research Intern',
            org: 'Rainmaker Capital',
            weight: 'mid',
            bullets: [
                'Built market and competitor research for the investment team.',
                'Turned findings into concise write-ups the team could act on.'
            ]
        },
        {
            period: 'Summer 2022',
            title: 'Server',
            org: 'Independent wine & beer bar',
            weight: 'early',
            summary: 'Guided guests through a rotating list of independent wines and beers.'
        },
        {
            period: '2021 — 2022',
            title: 'Lacrosse Coach',
            org: 'Youth lacrosse',
            weight: 'early',
            summary: 'Coached youth lacrosse, running practices and teaching fundamentals.'
        },
        {
            period: '2017 — 2021',
            title: 'Soccer Referee',
            org: 'Middle & high school',
            weight: 'early',
            summary: 'Officiated middle and high school matches, managing game flow and player safety.'
        }
    ];

    const present: Role[] = [
        {
            period: '2025 — Present',
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
            period: '2024 — Present',
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
            <div class="present">
                {#each present as role}
                    <article class="current-role">
                        <p class="current-period">{role.period}</p>
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

            <div class="fork-row" bind:this={forkRef}>
                <div class="fork-canvas" class:visible={forkVisible}>
                    <svg class="fork-desktop" viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true">
                        <path class="branch" d="M50 100 C50 48 24.2 62 24.2 0" pathLength="1" />
                        <path class="branch" d="M50 100 C50 48 75.8 62 75.8 0" pathLength="1" />
                    </svg>
                    <svg class="fork-mobile" viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true">
                        <path class="branch" d="M0 0 L0 100" pathLength="1" />
                    </svg>
                    <span class="junction" aria-hidden="true"></span>
                </div>
            </div>

            {#each history as role}
                <article
                    class="timeline-entry"
                    class:entry-recent={role.weight === 'recent'}
                    class:entry-mid={role.weight === 'mid'}
                    class:entry-early={role.weight === 'early'}
                >
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
        max-width: 1000px;
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

    .current-period {
        margin: 0 0 8px;
        font-family: var(--font-mono);
        font-size: 0.7rem;
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--text-muted);
    }

    .present {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 32px;
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
        bottom: -5px;
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
        font-size: 1.25rem;
    }

    .current-role .entry-org {
        font-size: 0.72rem;
    }

    .current-role .entry-summary,
    .current-role .entry-bullets li {
        font-size: 0.92rem;
    }

    .fork-row {
        position: relative;
        height: 110px;
    }

    .fork-canvas {
        position: absolute;
        inset: 0;
        clip-path: inset(0 0 100% 0);
        transition: clip-path 1s cubic-bezier(0.22, 1, 0.36, 1) 0.1s;
    }

    .fork-canvas.visible {
        clip-path: inset(0 0 0 0);
    }

    .fork-canvas svg {
        display: block;
        width: 100%;
        height: 100%;
        overflow: visible;
    }

    .fork-canvas .fork-desktop {
        display: block;
    }

    .fork-canvas .fork-mobile {
        display: none;
    }

    .branch {
        fill: none;
        stroke: var(--mint);
        stroke-width: 1;
        stroke-opacity: 0.55;
        vector-effect: non-scaling-stroke;
    }

    .junction {
        position: absolute;
        left: 50%;
        bottom: 0;
        width: 9px;
        height: 9px;
        border-radius: 50%;
        background: var(--mint);
        transform: translate(-50%, 50%);
    }

    .timeline-entry {
        --title-size: 1.12rem;
        --org-size: 0.7rem;
        --body-size: 0.9rem;
        --entry-pad: 32px;
        --node-size: 8px;
        --date-size: 0.72rem;
        position: relative;
        display: grid;
        grid-template-columns: 1fr 1fr;
    }

    .timeline-entry.entry-recent {
        --title-size: 1.35rem;
        --org-size: 0.75rem;
        --body-size: 0.95rem;
        --entry-pad: 44px;
        --node-size: 9px;
        --date-size: 0.75rem;
    }

    .timeline-entry.entry-mid {
        --title-size: 1.08rem;
        --org-size: 0.7rem;
        --body-size: 0.87rem;
        --entry-pad: 30px;
        --node-size: 8px;
        --date-size: 0.7rem;
    }

    .timeline-entry.entry-early {
        --title-size: 0.9rem;
        --org-size: 0.64rem;
        --body-size: 0.8rem;
        --entry-pad: 22px;
        --node-size: 6px;
        --date-size: 0.66rem;
    }

    .timeline-entry::before {
        content: '';
        position: absolute;
        left: 50%;
        top: 0;
        bottom: 0;
        width: 1px;
        background: var(--border);
    }

    .entry-date {
        grid-column: 1;
        text-align: right;
        padding: 3px 32px 0 0;
        font-family: var(--font-mono);
        font-size: var(--date-size);
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--text-muted);
        line-height: 1.5;
    }

    .node {
        position: absolute;
        left: 50%;
        top: 6px;
        width: var(--node-size);
        height: var(--node-size);
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
        padding: 0 0 var(--entry-pad) 32px;
    }

    .entry-title {
        margin: 0;
        font-family: var(--font-display);
        font-size: var(--title-size);
        font-weight: 600;
        color: var(--text);
        line-height: 1.3;
    }

    .entry-org {
        margin: 5px 0 0;
        font-family: var(--font-mono);
        font-size: var(--org-size);
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--mint);
    }

    .entry-summary {
        margin: 10px 0 0;
        color: var(--text-muted);
        font-size: var(--body-size);
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
        font-size: var(--body-size);
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

        .present {
            position: relative;
            grid-template-columns: 1fr;
            gap: 20px;
            padding-left: 22px;
        }

        .present::before {
            content: '';
            position: absolute;
            left: 0;
            top: 0;
            bottom: 0;
            width: 1px;
            background: var(--border);
        }

        .current-role::before {
            top: 26px;
            bottom: auto;
            left: -22px;
            transform: translate(-50%, -50%);
        }

        .fork-row {
            height: 56px;
        }

        .fork-row::before {
            content: '';
            position: absolute;
            left: 0;
            top: 0;
            bottom: 0;
            width: 1px;
            background: var(--border);
        }

        .fork-canvas .fork-desktop {
            display: none;
        }

        .fork-canvas .fork-mobile {
            display: block;
        }

        .fork-canvas {
            position: relative;
            height: 100%;
        }

        .junction {
            left: 0;
            bottom: auto;
            top: 0;
            transform: translate(-50%, -50%);
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
        }

        .entry-date {
            text-align: left;
            padding: 0 0 8px;
        }

        .entry-content {
            padding: 0 0 var(--entry-pad);
        }
    }
</style>
