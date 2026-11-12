<script lang="ts">
    import { onMount } from 'svelte';
    import { ChevronDown, ArrowUpRight } from 'lucide-svelte';

    type Weight = 'recent' | 'mid' | 'early';

    interface Role {
        id: string;
        period: string;
        title: string;
        org: string;
        weight?: Weight;
        blurb: string;
        bullets?: string[];
        tags?: string[];
        links?: { label: string; href: string }[];
    }

    // Placeholder dates and copy — replace with the real details.
    const history: Role[] = [
        {
            id: 'swe-intern',
            period: 'May 2026 — September 2026',
            title: 'Software Engineering Intern',
            org: 'Rainmaker Capital',
            weight: 'recent',
            blurb: 'Built internal systems and data-heavy features for the investment team.',
            bullets: [
                'REST API integration with the CRM, notes sheet generation, and automated email workflows.',
                'Developed a business valuation calculator that automated data intake. Fixed scoring edge cases and zero-value input validation bugs.',
                'Saved 1-2 hours of analyst time and worked across teams to deliver a better product.'
            ],
            tags: ['PHP', 'JavaScript', 'SQL', 'CRM']
        },
        {
            id: 'market-research',
            period: 'June 2024 - August 2024',
            title: 'Market Research Intern',
            org: 'Rainmaker Capital',
            weight: 'mid',
            blurb: 'Market research in the AI healthcare space.',
            bullets: [
                'Conducted indepth research on healthcare technologies and AI application.',
                'Turned research int a stakeholder report.'
            ]
        },
        {
            id: 'server',
            period: 'June 2024 - May 2026',
            title: 'Server and Storage Assistant',
            org: 'Indie wine & beer bar at SF Wine Center',
            weight: 'early',
            blurb: 'Served guests and managed inventory.',
            bullets: [
                'Serve food and beverages in a high-paced environment, and manage deliveries, while ensuring timely and accurate receipt of products',
                'Managed stock by organizing wine storage lockers and maintianing accurate inventory records.'
            ]
        },
        {
            id: 'lacrosse-coach',
            period: 'June 2023 — August 2023',
            title: 'Lacrosse Coach',
            org: 'Lacrosse the Bay',
            weight: 'early',
            blurb: 'Coached youth camps and club teams.',
            bullets: [
                'Attempted to tell kids what to do.',
                'Led and organized practice sessions and taught fundamentals to youth players.'
            ]
        },
        {
            id: 'referee',
            period: '2017 — 2022',
            title: 'Soccer Referee',
            org: 'Middle & high school',
            weight: 'early',
            blurb: 'Officiated youth soccer matches.',
            bullets: ['Got yelled at a lot.']
        }
    ];

    const present: Role[] = [
        {
            id: 'dataannotation',
            period: 'September 2026 — Present',
            title: 'Software Engineer',
            org: 'DataAnnotation (contract)',
            blurb: 'Contract software engineer evaluating and reviewing coding agents for AI training.',
            bullets: [
                'Evaluate AI coding agents on real-world software engineering tasks, including debugging large codebases, implementing features, and resolving complex technical issues.',
                'Identify failure modes in agentic workflows, including flawed reasoning, missed edge cases, and incorrect implementations.',
                'Review and assess agent generated code for correctness, code quality, maintainability, and adherence to requirements.'
            ],
            links: [{ label: 'dataannotation.tech', href: 'https://www.dataannotation.tech' }]
        },
        {
            id: 'metcode',
            period: 'August 2026 — Present',
            title: 'Founder & Maintainer',
            org: 'MetCode (open source)',
            blurb: 'Open-source agentic command-line assistant written in Go.',
            bullets: [
                'Created and maintain an open-source CLI agent optimized for small (7–8B parameter) models running offline on modest consumer hardware. 5+ outside contributors and 5+ stars on GitHub.',
                'Implemented the tool-calling loop including, tool dispatch, tool-result feedback, multi-call handling, and a file-system toolset. Test coverage at 60%+.',
                'Diagnosed model-specific failures: hallucinated tool names, and a model whose tool calls never parsed due to a template mismatch. Compared models on tool-call judgment and standardized on the best performer.'
            ],
            tags: ['Go', 'Ollama', 'Tailscale'],
            links: [{ label: 'GitHub', href: 'https://github.com/Edwin-metcalf/MetCode' }]
        }
    ];

    let openRoles: Record<string, boolean> = {};
    let headingVisible = false;
    let forkVisible = false;
    let sectionRef: HTMLElement | null = null;
    let forkRef: HTMLElement | null = null;

    function toggle(id: string) {
        openRoles = { ...openRoles, [id]: !openRoles[id] };
    }

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
                    {@const open = !!openRoles[role.id]}
                    <article class="current-role">
                        <button
                            class="entry-toggle"
                            type="button"
                            aria-expanded={open}
                            aria-controls={`detail-${role.id}`}
                            on:click={() => toggle(role.id)}
                        >
                            <span class="current-period">{role.period}</span>
                            <span class="entry-head">
                                <span class="entry-title">{role.title}</span>
                                <span class="chevron" aria-hidden="true"><ChevronDown size={16} /></span>
                            </span>
                            <span class="entry-org">{role.org}</span>
                            <span class="entry-blurb">{role.blurb}</span>
                        </button>

                        <div class="entry-detail" id={`detail-${role.id}`} inert={!open}>
                            <div class="entry-detail-inner">
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
                                {#if role.links?.length}
                                    <div class="entry-links">
                                        {#each role.links ?? [] as link}
                                            <a
                                                class="entry-link"
                                                href={link.href}
                                                target="_blank"
                                                rel="noreferrer"
                                            >
                                                {link.label}
                                                <ArrowUpRight size={13} />
                                            </a>
                                        {/each}
                                    </div>
                                {/if}
                            </div>
                        </div>
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
                </div>
            </div>

            {#each history as role}
                {@const open = !!openRoles[role.id]}
                <article
                    class="timeline-entry"
                    class:entry-recent={role.weight === 'recent'}
                    class:entry-mid={role.weight === 'mid'}
                    class:entry-early={role.weight === 'early'}
                >
                    <span class="node" aria-hidden="true"></span>
                    <div class="entry-date">{role.period}</div>
                    <div class="entry-content">
                        <button
                            class="entry-toggle"
                            type="button"
                            aria-expanded={open}
                            aria-controls={`detail-${role.id}`}
                            on:click={() => toggle(role.id)}
                        >
                            <span class="entry-head">
                                <span class="entry-title">{role.title}</span>
                                <span class="chevron" aria-hidden="true"><ChevronDown size={16} /></span>
                            </span>
                            <span class="entry-org">{role.org}</span>
                            <span class="entry-blurb">{role.blurb}</span>
                        </button>

                        <div class="entry-detail" id={`detail-${role.id}`} inert={!open}>
                            <div class="entry-detail-inner">
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
                                {#if role.links?.length}
                                    <div class="entry-links">
                                        {#each role.links ?? [] as link}
                                            <a
                                                class="entry-link"
                                                href={link.href}
                                                target="_blank"
                                                rel="noreferrer"
                                            >
                                                {link.label}
                                                <ArrowUpRight size={13} />
                                            </a>
                                        {/each}
                                    </div>
                                {/if}
                            </div>
                        </div>
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

    /* current roles: plain timeline entries in two columns */

    .present {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 32px;
    }

    .current-role {
        --title-size: 1.35rem;
        --org-size: 0.75rem;
        --body-size: 0.95rem;
        --node-size: 9px;
        position: relative;
        padding-bottom: 30px;
    }

    .current-role::before {
        content: '';
        position: absolute;
        bottom: -5px;
        left: 50%;
        width: var(--node-size);
        height: var(--node-size);
        border-radius: 50%;
        background: var(--mint);
        transform: translateX(-50%);
        z-index: 1;
    }

    .current-period {
        display: block;
        margin-bottom: 8px;
        font-family: var(--font-mono);
        font-size: 0.7rem;
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--text-muted);
    }

    /* clickable entry */

    .entry-toggle {
        display: block;
        width: 100%;
        padding: 0;
        background: none;
        border: none;
        color: inherit;
        font: inherit;
        text-align: left;
        cursor: pointer;
    }

    .entry-toggle:hover .entry-title,
    .entry-toggle:hover .chevron {
        color: var(--mint);
    }

    .entry-head {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 16px;
    }

    .entry-title {
        font-family: var(--font-display);
        font-size: var(--title-size);
        font-weight: 600;
        color: var(--text);
        line-height: 1.3;
        transition: color 0.15s ease;
    }

    .entry-org {
        display: block;
        margin-top: 5px;
        font-family: var(--font-mono);
        font-size: var(--org-size);
        letter-spacing: 0.5px;
        text-transform: uppercase;
        color: var(--mint);
    }

    .entry-blurb {
        display: block;
        margin-top: 8px;
        color: var(--text-muted);
        font-size: var(--body-size);
        line-height: 1.6;
        max-width: 62ch;
    }

    .chevron {
        flex: none;
        display: inline-flex;
        margin-top: 3px;
        color: var(--text-muted);
        transition: transform 0.25s ease, color 0.15s ease;
    }

    .entry-toggle[aria-expanded='true'] .chevron {
        transform: rotate(180deg);
    }

    /* expandable detail */

    .entry-detail {
        display: grid;
        grid-template-rows: 0fr;
        transition: grid-template-rows 0.3s cubic-bezier(0.22, 1, 0.36, 1);
    }

    .entry-detail-inner {
        overflow: hidden;
        min-height: 0;
    }

    .entry-toggle[aria-expanded='true'] + .entry-detail {
        grid-template-rows: 1fr;
    }

    .entry-bullets {
        list-style: none;
        margin: 14px 0 0;
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

    .entry-links {
        display: flex;
        flex-wrap: wrap;
        gap: 16px;
        margin-top: 14px;
    }

    .entry-link {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-family: var(--font-mono);
        font-size: 0.72rem;
        letter-spacing: 0.3px;
        color: var(--mint);
        text-decoration: none;
        border-bottom: 1px solid transparent;
        transition: border-color 0.15s ease;
    }

    .entry-link:hover {
        border-bottom-color: var(--mint);
    }

    /* fork */

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

    /* history entries */

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

        .current-role {
            padding-bottom: 24px;
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
