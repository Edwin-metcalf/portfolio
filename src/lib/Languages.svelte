<script lang="ts">
    import {
        siJavascript,
        siOpenjdk,
        siPython, 
        siTypescript,
        siOcaml,
        siC,
        siRuby,
        siGo,
        siHtml5,
        siCss,
        siGit,
        siHeroku,
        siRubyonrails,
        siGodotengine,
        siSvelte
    } from 'simple-icons';
    const skills = [
    // Languages
    { icon: siPython, name: 'Python', color: '#3776AB' },
    { icon: siJavascript, name: 'JavaScript', color: '#F7DF1E' },
    { icon: siTypescript, name: 'TypeScript', color: '#3178C6' },
    { icon: siOcaml, name: 'OCaml', color: '#EC6813' },
    { icon: siC, name: 'C', color: '#A8B9CC' },
    { icon: siRuby, name: 'Ruby', color: '#CC342D' },
    { icon: siGo, name: 'Go', color: '#00ADD8' },
    { icon: siOpenjdk, name: 'Java', color: '#007396' },
    { icon: siHtml5, name: 'HTML5', color: '#E34F26' },
    { icon: siCss, name: 'CSS3', color: '#1572B6' },
    
    // Tools
    { icon: siGit, name: 'Git', color: '#F05032' },
    { icon: siHeroku, name: 'Heroku', color: '#430098' },
    { icon: siRubyonrails, name: 'Rails', color: '#CC0000' },
    { icon: siGodotengine, name: 'Godot', color: '#478CBF' },
    { icon: siSvelte, name: 'Svelte', color: '#FF3E00'}
];

</script>

<div class="skills-sections" id="skills">
    <div class="skills-container">
        <h1 class="skills-title">Languages and Tools</h1>

        <div class="skills-category">
            <div class="skills-grid">
                {#each skills as lang}
                    <div class="skill-item">
                        <div class="skill-icon" style="color: {lang.color}">
                            {@html lang.icon.svg}
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>

<style>
    .skills-sections {
        padding: 80px 5%;
        display: flex;
        justify-content: center;
    }
    .skills-container {
        max-width: 1000px;
        width: 100%;
    }
    .skills-title {
        font-size: 3rem;
        font-weight: 600;
        color: #f0f6fc;
        text-align: center;
        margin-bottom: 50px;
        position: relative;
    }
    .skills-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 3px;
        background: #00d4aa;
        margin: 20px auto;
        border-radius: 2px;
    }
    .skills-category {
        margin-bottom: 50px;
    }
    .skills-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
        gap: 30px;
        justify-items: center;
    }

    .skill-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 20px;
        background: rgba(33, 38, 45, 0.5);
        border-radius: 12px;
        transition: all 0.3s ease;
        border: 1px solid rgba(48, 54, 61, 0.5);
        min-width: 120px;
    }

    .skill-item:hover {
        transform: translateY(-5px);
        background: rgba(33, 38, 45, 0.8);
        border-color: #00d4aa;
        box-shadow: 0 10px 25px rgba(0, 212, 170, 0.1);
    }

    .skill-icon {
        width: 48px;
        height: 48px;
        margin-bottom: 15px;
        transition: transform 0.3s ease;
    }

    .skill-item:hover .skill-icon {
        transform: scale(1.1);
    }

    .skill-icon :global(svg) {
        width: 100%;
        height: 100%;
        fill: currentColor;
    }
    
</style>