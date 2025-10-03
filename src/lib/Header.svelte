<script lang="ts">
</script>

<div class="header">
    <a href="/" class="logo-wrapper">
        <img src="/images/EMlogo.png" alt="Edwin Metcalf Logo" class="logo" />
    </a>
    <div class="header-right">
        <a href="#about-me">About me</a>
        <a href="#skills">Skills</a>
        <a href="#portfolio">Portfolio</a>
        <a href="#contact">Contact</a>
    </div>

</div>

<style>
    .header {
        overflow: hidden;
        background: linear-gradient(90deg, #0d1117, #161b22, #21262d);
        padding: 20px 10px;
        color: #f0f6fc;
        transform: translateY(-2px);
    }
    .header a {
        float: left;
        color: #f0f6fc;
        text-align: center;
        padding: 12px;
        text-decoration: none;
        font-size: 18px;
        line-height: 25px;
        border-radius: 4px;
    }
    .header a:hover {
        background-color: #00d4aa;
        color: #0d1117;
        transform: translateY(-2px);
    }
    .header-right {
        float: right;
    }
    .logo {
        width: 50px;
        height: auto;
        display: block;
        transition: filter 0.3s ease;
    }
    .logo:hover {
        filter: brightness(0.5);
    }
</style>