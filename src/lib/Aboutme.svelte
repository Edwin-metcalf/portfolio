<script lang="ts">
</script>
<div class="about-section" id="about-me">
    <div class="about-container">
    <h1 class="about-title">About Me</h1>
        <p class="about-text">
            Hi welcome to my website. I am an aspiring full stack developer. I am currently am studying both computer science and history. I am from San Francisco and have always been super interested in computers and especially video games.
        </p>
    </div>
</div>

<style>
    .about-section {
        padding: 80px 5%;
        display: flex;
        justify-content: center;
    }
    .about-container {
        text-align: center;
        max-width: 800px;
    }
    .about-title {
        font-size: 3rem;
        font-weight: 600;
        color: #f0f6fc;
        margin-bottom: 30px;
        position: relative;
    }
    .about-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 3px;
        background: #00d4aa;
        margin: 20px auto;
        border-radius: 2px;
    }
    .about-text {
        font-size: 1.2rem;
        line-height: 1.8;
        color: #c9d1d9;
        margin: 0;
    }
</style>