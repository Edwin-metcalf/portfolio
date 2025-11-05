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

    private testEnemy: Enemy;

    constructor(canvas: HTMLCanvasElement) {
        this.canvas = canvas;
        const context = canvas.getContext("2d");

        if (!context) throw new Error("could not get context");
        this.c = context;

        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;

        this.player = new Player(this.c, this.canvas.width, this.canvas.height);
        this.testEnemy = new Enemy(this.c);
        this.testEnemy.draw();
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
                this.projectiles.push(new Projectile(this.c, {x: this.player.position.x + this.player.width / 2 ,y: this.player.position.y + 30 }, {x: 0 ,y: -3  }));


            }
        }

    }

    private spawnEnemyGrid(): void {
        const rows = 4;
        const cols = 15;
        const enemySpacing = 60;
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
                        }

                    }, 0);
                }
            });
            enemy.update();
        });
    }
    
    animate = (): void => {
        this.animationId = requestAnimationFrame(this.animate);
        this.c.fillStyle = 'black';
        this.c.fillRect(0,0, this.canvas.width, this.canvas.height);
        this.updatePlayerControls();
        this.player.update();
        //this.testEnemy.update();
        this.updateEnemies();

        this.projectiles.forEach((projectile, idx) => {
            if (projectile.position.y + projectile.radius <= 0 ) {
                setTimeout(() => {
                    this.projectiles.splice(idx, 1);

                }, 0);
            } else {
                projectile.update()
            }
            
        });


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

class Projectile {
    public position: Position;
    public velocity: Velocity;
    public radius: number;
    private c: CanvasRenderingContext2D;

    constructor(context: CanvasRenderingContext2D, position: Position, velocity: Velocity ) {
        this.position = position;
        this.velocity = velocity;
        this.c = context;

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
        this.position.x += this.velocity.x;
        this.position.y += this.velocity.y;  
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

        this.width = image.width * 0.10;
        this.height = image.height * 0.10;
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