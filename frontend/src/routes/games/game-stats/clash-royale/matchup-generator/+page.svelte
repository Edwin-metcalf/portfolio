<script lang="ts">
	import { fetchAPI } from "$lib/api";
    import { type ClashRoyaleFriendlyLoadReturn, type matchupGeneratorTags } from "../../stats";

    let matchupEntry = $state<matchupGeneratorTags>({tag1: "", tag2: ""});
    let tagsSubmitted: boolean = $state(false)



    async function getMatchUpData(tags: matchupGeneratorTags): Promise<ClashRoyaleFriendlyLoadReturn | null> {
        try {
            const result = await fetchAPI('/api.games/clash-royale/matchup-generator', { 
                method: 'PUT',
                body: JSON.stringify(tags)

        });
        return result as ClashRoyaleFriendlyLoadReturn;
        } catch (err) {
            console.error("Error: ", err);
            return null;
        }

    }

</script>
<div>
    <h1>
        Generate your matchup against you and your friend!
    </h1>
    <div class="tag-submisssion">
        <input type="text" bind:value={matchupEntry.tag1} placeholder="Your Clash Royale Tag" class="tag-input"/>
        <input type="text" bind:value={matchupEntry.tag2} placeholder="Opponents Clash Royale Tag" class="tag-input"/>

        {#if !tagsSubmitted}
            <button class="submit-button" onclick={() => getMatchUpData(matchupEntry)}>Submit </button>
        {/if}

    </div>
</div>