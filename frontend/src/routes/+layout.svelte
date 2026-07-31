<script lang="ts">
    import '../lib/styles/theme.css';


	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href="/images/EMlogo.png"/>
</svelte:head>

{@render children?.()}
