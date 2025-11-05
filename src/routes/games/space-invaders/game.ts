import { Play } from "lucide-svelte";
import playerImageUrl from './imgs/player.png';
import { siVelocity } from "simple-icons";
import enemyImg from './imgs/enemy.png';

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

    public score: number;
    private onScoreChange?: (score: number) => void;

    private player: Player;
    private animationId: number | null = null;
    private keys: { [key: string]: boolean} = {};

    //variables for the shooting 
    private projectiles: Projectile[] = [];
    private lastShotTime: number = 0;
    private shootCooldown: number = 250;

    //variables for the enemies
    private enemies: Enemy[] = [];
    private enemyVelocity: Velocity = {x:3, y:0};
    private enemyDropDistance: number = 40;

    private particles: Particle[] = [];
    private game = {
        over: false,
        active: false
    }


    constructor(canvas: HTMLCanvasElement, onScoreChange?: (score: number) => void) {
        this.canvas = canvas;
        const context = canvas.getContext("2d");
        this.onScoreChange = onScoreChange;
        this.score = 0;

        if (!context) throw new Error("could not get context");
        this.c = context;

        this.canvas.width = 1024;
        this.canvas.height = 576;

        this.player = new Player(this.c, this.canvas.width, this.canvas.height);
    
        this.player.draw();
        this.spawnEnemyGrid();
        
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

    private updatePlayerControls(): void {
        if (!(this.game.over)) {
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
            if(this.keys[' ']) {
                const currentTime = Date.now();
                if (currentTime - this.lastShotTime >= this.shootCooldown){
                    this.lastShotTime = currentTime;
                    this.projectiles.push(new Projectile(this.c, {x: this.player.position.x + this.player.width / 2 ,y: this.player.position.y + 30 }, {x: 0 ,y: -1  }));


                }
            }
        }   

    }

    private spawnEnemyGrid(): void {
        const rows = 4;
        const cols = 15;
        const enemySpacing = 30;
        const startX = 100;
        const startY = 50;

        for (let row = 0; row < rows; row++){
            for (let col = 0; col < cols; col++) {
                this.enemies.push(new Enemy(this.c, {x: startX + col * enemySpacing, y: startY + row * enemySpacing}));
            }
        }


    }

    private updateEnemies(): void {
        let hitEdge = false;

        this.enemies.forEach(enemy => {
            if ((enemy.position.x + enemy.width >= this.canvas.width && this.enemyVelocity.x > 0) || (enemy.position.x <= 0 && this.enemyVelocity.x < 0)) {
                hitEdge = true;
            }
        });

        if (hitEdge) {
            this.enemyVelocity.x *= -1;
            this.enemies.forEach(enemy => {
                enemy.position.y += this.enemyDropDistance;
            })
        }
        this.enemies.forEach((enemy, i) => {
            enemy.velocity = {...this.enemyVelocity};
            this.projectiles.forEach((projectile, j) => {
                if (
                    projectile.position.y - projectile.radius <= enemy.position.y + enemy.height && 
                    projectile.position.x + projectile.radius >= enemy.position.x &&
                    projectile.position.x - projectile.radius <= enemy.position.x + enemy.width &&
                    projectile.position.y + projectile.radius >= enemy.position.y
                ) {
                    setTimeout(() => {
                        const invaderFound = this.enemies.find(enemy2 => enemy2 === enemy );

                        const projectileFound = this.projectiles.find(projectile2 => projectile2 === projectile);

                        if (invaderFound && projectileFound){
                            this.enemies.splice(i, 1);
                            this.projectiles.splice(j,1);

                            this.createParticles(enemy);
                            this.score += 100;
                            if (this.onScoreChange) {
                                this.onScoreChange(this.score)
                            }
                        } 

                    }, 0);
                }
            });
            if (
                enemy.position.x < this.player.position.x + this.player.width &&
                enemy.position.x + enemy.width  > this.player.position.x  &&
                enemy.position.y  < this.player.position.y &&
                enemy.position.y + enemy.height > this.player.position.y
                ) {

                    console.log("game over type beat");
                    this.player.isAlive = false;
                    if (!(this.game.over)){
                        for (let i = 0; i < 15; i++) {
                            this.createParticles(this.player, 'grey');
                        }
                    }
                    this.stop()
                }
            enemy.update();
        });
    }

    createParticles(object: Enemy | Player , color: string = 'purple'){
        for (let i = 0; i < 15; i++) {
            this.particles.push(new Particle(
            this.c,
            {x: object.position.x + object.width / 2, y: object.position.y + object.height / 2}, 
            { x: (Math.random() - 0.5) * 2, y: (Math.random() - 0.5) * 2},
            Math.random() * 3,
            color
            ));
        }

    }
    
    animate = (): void => {
        this.animationId = requestAnimationFrame(this.animate);
        this.c.fillStyle = 'black';
        this.c.fillRect(0,0, this.canvas.width, this.canvas.height);
        this.updatePlayerControls();
        if (this.game.over) {
            setTimeout(() => {
                return;
            }, 40)
        } else {
            this.player.update();

        }
        this.updateEnemies();

        this.particles.forEach((particle, i) => {
            if (particle.opacity <= 0) {
                setTimeout(() => {
                    this.particles.splice(i,1);

                }, 0);
            } else {
                particle.update();
            }
        });

        this.projectiles.forEach((projectile, idx) => {
            if (projectile.position.y + projectile.radius <= 0 ) {
                setTimeout(() => {
                    this.projectiles.splice(idx, 1);

                }, 0);
            } else {
                projectile.update()
            }
            
        });     

        for (let i = 0; i < 1; i++) {
            this.particles.push(new Particle(
            this.c,
            {x: Math.random() * this.canvas.width, y: Math.random() * this.canvas.height}, 
            { x: 0, y: 1},
            Math.random() * 3,
            'white'
            ));
        
    }


    }
    start(): void {
        this.game.over = false;
        this.animate();
    }
    stop(): void {
        this.game.over = true;

        setTimeout(() => {

        
            if (this.animationId != null) {
                cancelAnimationFrame(this.animationId);
                this.animationId = null;
            }

            window.removeEventListener('keydown', this.handleKeyDown);
            window.removeEventListener('keyup', this.handleKeyUp);
        }, 2000);
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
    public isAlive: boolean;


    constructor(context: CanvasRenderingContext2D, canvasWidth: number, canvasHeight: number) {
        this.c = context;
        this.rotation = 0;
        this.isAlive = true;

        const image = new Image();
        image.src = playerImageUrl
        image.onload = () => {
            //this.imageLoaded = true;
            console.log('Player image loaded!');
        };
        this.image = image;
        this.width = image.width * 0.1;
        this.height = image.height * 0.1;

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

class Projectile {
    public position: Position;
    public velocity: Velocity;
    public radius: number;
    private c: CanvasRenderingContext2D;
    private speed: number;

    constructor(context: CanvasRenderingContext2D, position: Position, velocity: Velocity ) {
        this.position = position;
        this.velocity = velocity;
        this.c = context;
        
        this.speed = 5;
        this.radius = 3;

    }
    draw() {
        this.c.beginPath();
        this.c.arc(this.position.x, this.position.y, this.radius, 0, Math.PI * 2);

        this.c.fillStyle = 'red';
        this.c.fill();
        this.c.closePath();
    }

    update() {
        this.draw()
        this.position.x += this.velocity.x * this.speed;
        this.position.y += this.velocity.y * this.speed;  
    }

}

class Particle {
    public position: Position;
    public velocity: Velocity;
    public radius: number;
    private c: CanvasRenderingContext2D;
    public color: string;
    public opacity: number;

    constructor(context: CanvasRenderingContext2D, position: Position, velocity: Velocity, radius: number, color: string ) {
        this.position = position;
        this.velocity = velocity;
        this.c = context;

        this.radius = radius;
        this.color = color;

        this.opacity = 1;
    }
    draw() {
        this.c.save()
        this.c.globalAlpha = this.opacity;
        this.c.beginPath();
        this.c.arc(this.position.x, this.position.y, this.radius, 0, Math.PI * 2);

        this.c.fillStyle = this.color;
        this.c.fill();
        this.c.closePath();
        this.c.restore()
    }

    update() {
        this.draw()
        this.position.x += this.velocity.x;
        this.position.y += this.velocity.y; 
        this.opacity -= 0.01; 
    }
}

class Enemy {
    private c: CanvasRenderingContext2D;
    public position: Position;
    public velocity: Velocity;
    private image: HTMLImageElement;
    public readonly width: number;
    public readonly height: number;

    constructor(context: CanvasRenderingContext2D, startPosition?: Position) {
        this.c = context;
        const image = new Image();
        image.src = enemyImg;

        this.image = image;

        this.width = image.width * 0.05;
        this.height = image.height * 0.05;
        this.velocity = { x: 0, y: 0 };
        this.position = startPosition || { x: 300, y: 300};

    }
    draw(): void {
        this.c.save();
        this.c.filter = 'hue-rotate(300deg) saturate(2) brightness(1.1)';
        this.c.drawImage(this.image, this.position.x, this.position.y, this.width, this.height);
        this.c.restore();
    }
    update(): void {
        this.draw()
        this.position.x += this.velocity.x;
        this.position.y += this.velocity.y;
    }
}