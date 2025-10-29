import { Play } from "lucide-svelte";
import playerImageUrl from './imgs/player.png';

interface Position {
    x: number;
    y: number;
}

interface Velocity {
    x: number;
    y: number;
}
export class SpaceInvaders {
    private canvas: HTMLCanvasElement;
    private c: CanvasRenderingContext2D;
    private player: Player;
    private animationId: number | null = null;
    private keys: { [key: string]: boolean} = {};

    constructor(canvas: HTMLCanvasElement) {
        this.canvas = canvas;
        const context = canvas.getContext("2d");

        if (!context) throw new Error("could not get context");
        this.c = context;

        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;

        this.player = new Player(this.c, this.canvas.width, this.canvas.height);
        this.player.draw();
        
        this.setupEventListeners();
    }

    private setupEventListeners(): void {
        window.addEventListener('keydown', this.handleKeyDown);
        window.addEventListener('keyup', this.handleKeyUp);
    }
    private handleKeyDown = (event: KeyboardEvent): void => {
        this.keys[event.key] = true;
    }
    private handleKeyUp = (event: KeyboardEvent): void => {
        this.keys[event.key] = false;
    }
    private updatePlayerMovement(): void {
        this.player.rotation = 0;
        this.player.velocity.x = 0;
        const speed = 5;
        if (this.keys['ArrowLeft'] || this.keys['a']) {
            this.player.velocity.x = -speed;
            this.player.rotation = -.2;
        }
        if (this.keys['ArrowRight'] || this.keys['d']) {
            this.player.velocity.x = speed;
            this.player.rotation = .2;
        }

    }
    
    animate = (): void => {
        this.animationId = requestAnimationFrame(this.animate);
        this.c.fillStyle = 'black';
        this.c.fillRect(0,0, this.canvas.width, this.canvas.height);
        this.updatePlayerMovement();
        this.player.update();


    }
    start(): void {
        this.animate();
    }
    stop(): void {
        if (this.animationId != null) {
            cancelAnimationFrame(this.animationId);
            this.animationId = null;
        }

        window.removeEventListener('keydown', this.handleKeyDown);
        window.removeEventListener('keyup', this.handleKeyUp);
    }

}

class Player {
    public position: Position;
    public velocity: Velocity;
    public readonly width: number;
    public readonly height: number;
    private c: CanvasRenderingContext2D;
    private image: HTMLImageElement;
    public rotation: number;


    constructor(context: CanvasRenderingContext2D, canvasWidth: number, canvasHeight: number) {
        this.c = context;
        this.rotation = 0;

        const image = new Image();
        image.src = playerImageUrl
        image.onload = () => {
            //this.imageLoaded = true;
            console.log('Player image loaded!');
        };
        this.image = image;
        this.width = image.width * 0.15;
        this.height = image.height * 0.15;

        this.position = { x: canvasWidth / 2 - this.width / 2, y: canvasHeight - this.height - 20 };
        this.velocity = { x: 0, y: 0 };
    }

    draw(): void {
        this.c.save()
        this.c.translate(this.position.x + this.width / 2, this.position.y + this.height / 2);
        this.c.rotate(this.rotation)
        this.c.translate(-this.position.x - this.width / 2, -this.position.y - this.height / 2);


        if (this.image) this.c.drawImage(this.image, this.position.x, this.position.y, this.width, this.height);

        this.c.restore()
    }
    update(): void {
        this.position.x += this.velocity.x;
        this.position.y += this.velocity.y;   
        this.draw();
    }
    
}