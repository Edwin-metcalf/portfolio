<script lang="ts">
    import { X } from "lucide-svelte";

    const games = [
        {
            title: "Space Invaders",
            description: "My take on the classic space invaders",
            url: "/games/space-invaders",
            status: "playable"
        },
        {
            title: "Under Construction",
            description: "always working on more cool stuff",
            url: "#",
            status: "coming-soon"
        }
    ];



</script>
<!-- what the fuck was that bro it changed the whole home screen as well???
<div class="game-header">
    <a href="/" class="logo-link">
        <img src="/images/EMlogo.png" alt="logo navigate to homepage" class="logo"/>
    </a>
</div> -->

<div class="games-section"> 
    <a href="/" class="exit-button" aria-label="Exit games">
        <X size={24} />
    </a>
    <h1>Games & Projects</h1>

    <div class="games-grid">
        {#each games as game}
            <a
                href={game.url}
                class="game-card"
                class:coming-soon={game.status === "comingssoon"}
            >
                <div class="game-content">
                    <h2>{game.title}</h2>
                    <p>{game.description}</p>
                    {#if game.status === "playable"}
                        <span class="play-button">Play Now</span>
                    {:else}
                        <span class="coming-soon">Coming soon</span>
                    {/if}
                </div>
            </a>
        {/each}
    </div>
</div>

<style>
    .games-section {
        position: relative;
        margin: 0 auto;
        padding: 4rem 2rem;
        min-height: 100vh;
        background-color: #161b22;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }
    h1 {
        font-size: 3rem;
        margin-bottom: 2rem;
        text-align: center;
        color: #00d4aa;
        width: 100%;
    }
    .games-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 2rem;
        justify-content: center;
    }
    .game-card {
        background: linear-gradient(135deg, #161b22 0%, #21262d 100%);;
        border-radius: 12px;
        padding: 2rem;
        text-decoration: none;
        color: white;
        transition: transform 0.3s ease, box-shadow 0.3s ease;
        cursor: pointer;
        position: relative;
        overflow: hidden;
        border: 1px solid rgba(0, 212, 170, 0.3);
    }
    .game-card:hover {
        transform: translateY(-5px);
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
    }
    .game-card.coming-soon {
        background: linear-gradient(135deg, #434343 0%, #000000 100%);
        cursor: not-allowed;
        opacity: 0.6;
    }
    .game-card.coming-soon:hover {
        transform: none;
    }
    .game-content h2 {
        font-size: 1.8rem;
        margin-bottom: 0.5rem;
    }

    .game-content p {
        font-size: 1rem;
        margin-bottom: 1rem;
        opacity: 0.9;
    }
    .play-button {
        display: inline-block;
        background: rgba(0, 212, 170, 0.15);
        padding: 0.5rem 1rem;
        border-radius: 6px;
        font-weight: bold;
        color: #00d4aa;
    }

    .game-card:hover .play-button {
        background: rgba(0, 212, 170, 0.25); 
    }

    .coming-soon {
        display: inline-block;
        background: rgba(255, 255, 255, 0.1);
        padding: 0.5rem 1rem;
        border-radius: 6px;
        font-size: 0.9rem;
    }

    .exit-button {
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
    
    .exit-button:hover {
        background: rgba(0, 212, 170, 0.2);
        border-color: #00d4aa;
        transform: rotate(90deg);

    }


    @media (max-width: 768px) {
        h1 {
            font-size: 2rem;
        }
        
        .games-grid {
            grid-template-columns: 1fr;
        }
    }

</style>
