<script lang="ts">
    import { Mail, Github, Linkedin, Send } from 'lucide-svelte';
    let formData = {
        name: '',
        email: '',
        message: ''
    };

    let isSubmitting = false;
    let submitStatus = '';

    async function handleSubmit(event: Event) {
        event.preventDefault();
        isSubmitting = true;

        try {
            const response = await fetch('https://formspree.io/f/xovkevbp', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });
            if (response.ok) {
                submitStatus = 'success';
                formData = {name: '', email: '', message: ''};

            } else {
                submitStatus = 'error';
            }
        } catch (error) {
            submitStatus = 'error';
        }
        isSubmitting = false;
        setTimeout(() => submitStatus = '', 3000);
    }
</script>

<div class="contact-section" id="contact">
    <div class="contact-container">
        <div class="contact-header">
            <h1 class="contact-title">Get in Touch With Me</h1>
        </div>

        <div class="contact-content">
            <div class="contact-form-container">
                <form class="contact-form" on:submit={handleSubmit}>
                    <div class="form-group">
                        <label for="name">Name</label>
                        <input type="text" id="name" bind:value={formData.name} placeholder="Your name" required/>
                    </div>

                    <div class="form-group">
                        <label for="email">Email</label>
                        <input type="email" id="email" bind:value={formData.email} placeholder="your.email@example.com" required/>
                    </div>

                    <div class="form-group">
                        <label for="message">Message</label>
                        <textarea id="message" bind:value={formData.message} placeholder="Your message..." rows="5" required></textarea>
                    </div>
                    <button type="submit" class="submit-button" disabled={isSubmitting}>
                        {#if isSubmitting}
                            Sending...
                        {:else}
                            <Send size={20}/>
                            Send Message
                        {/if}
                    </button>
                    {#if submitStatus === 'success'}
                        <p class="status-message success">Message sent succesfully</p>
                    {:else if submitStatus === 'error'}
                        <p class="status-message error">Failed to send message please try again</p>
                    {/if}

                </form>
            </div>
            <div class="contact-info">
                <h3 class="info-title">Let's Connect</h3>
            
                <div class="contact-links">
                    <a href="mailto:winmetcalf1@gmail.com" class="contact-link">
                            <Mail size={20} />
                            winmetcalf1@gmail.com
                    </a>

                    <a href="https://github.com/Edwin-metcalf" class="contact-link" target="_blank">
                            <Github size={20} />
                            github.com/Edwin-metcalf
                    </a>
                    
                    <a href="https://www.linkedin.com/in/edwin-metcalf" class="contact-link" target="_blank">
                            <Linkedin size={20} />
                            linkedin.com/in/edwin-metcalf
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .contact-section {
        padding: 80px 5%;
        display: flex;
        justify-content: center;
    }

    .contact-container {
        max-width: 1200px;
        width: 100%;
    }

    .contact-header {
        text-align: center;
        margin-bottom: 60px;
    }

    .contact-title {
        font-size: 3rem;
        font-weight: 600;
        color: var(--text);
        text-align: center;
        margin-bottom: 50px;
        position: relative;
        font-family: var(--font-display);
    }
    
    .contact-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 1px;
        background: var(--border);
        margin: 20px auto;
    }

    .contact-content {
        display: grid;
        grid-template-columns: 2.3fr 1fr;
        gap: 60px;
        align-items: start;
    }

    .contact-form-container {
        background: var(--panel);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 40px;
    }

    .contact-form {
        display: flex;
        flex-direction: column;
        gap: 25px;
    }
    
    .form-group {
        display: flex;
        flex-direction: column;
        gap: 8px;
    }
    
    .form-group label {
        color: var(--text);
        font-weight: 500;
        font-size: 0.95rem;
        font-family: var(--font-mono);
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .form-group input,
    .form-group textarea {
        background: transparent;
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 15px;
        color: var(--text);
        font-size: 1rem;
        transition: all 0.15s ease;
        resize: vertical;
    }

    .form-group input::placeholder,
    .form-group textarea::placeholder {
        color: var(--text-muted);
    }
    
    .form-group input:focus,
    .form-group textarea:focus {
        outline: none;
        border-color: var(--mint);
    }

    .submit-button {
        background: transparent;
        color: var(--mint);
        border: 1px solid var(--mint);
        border-radius: var(--radius);
        padding: 15px 30px;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.15s ease;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;
        font-family: var(--font-mono);
    }
    
    .submit-button:hover:not(:disabled) {
        transform: translateY(-1px);
        background: rgba(62, 207, 192, 0.15);
    }
    .submit-button:disabled {
        opacity: 0.7;
        cursor: not-allowed;
    }
    
    .status-message {
        text-align: center;
        font-weight: 500;
        padding: 10px;
        border-radius: 8px;
        margin: 0;
    }
    
    .status-message.success {
        color: var(--mint);
        background: rgba(62, 207, 192, 0.15);
        border: 1px solid var(--mint);
    }
    
    .status-message.error {
        color: #ff6b6b;
        background: rgba(255, 107, 107, 0.1);
        border: 1px solid rgba(255, 107, 107, 0.3);
    }

    .contact-info {
        padding: 20px 0;
    }

    .info-title {
        font-size: 1.5rem;
        font-weight: 600;
        color: var(--text);
        margin-bottom: 20px;
        font-family: var(--font-display);
    }

    .contact-links {
        display: flex;
        flex-direction: column;
        gap: 15px;
        margin-bottom: 30px;
    }
    
    .contact-link {
        display: flex;
        align-items: center;
        gap: 12px;
        color: var(--text-muted);
        text-decoration: none;
        font-size: 1rem;
        transition: all 0.15s ease;
        padding: 8px 0;
    }
    
    .contact-link:hover {
        color: var(--mint);
    }

      @media (max-width: 768px) {
        .contact-section {
            padding: 80px 5%;
        }
        .contact-content {
            grid-template-columns: 1fr;
            gap: 40px;
        }
        
        .contact-form-container {
            padding: 30px 20px;
        }
        
        .contact-title {
            font-size: 2.5rem;
        }
        .contact-section {
            padding-left: 5%;
        }
    }

</style>