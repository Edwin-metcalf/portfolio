import { Play } from "lucide-svelte";
import playerImageUrl from './imgs/player.png';
import { siVelocity } from "simple-icons";
import enemyImg from './imgs/enemy.png';
import hunterImg from './imgs/enemyHunter.png'

interface Position {
    x: number;
    y: number;
}

interface Velocity {
    x: number;
    y: number;
}
type EnemyType = Enemy | HunterEnemy

export class SpaceInvaders {
    private canvas: HTMLCanvasElement;
    private c: CanvasRenderingContext2D;

    public score: number;
    private onScoreChange?: (score: number) => void;
    private onGameOver?: () => void;
    private onMessage: (message: string) => void;

    private player!: Player;
    private animationId: number | null = null;
    private keys: { [key: string]: boolean} = {};

    //variables for the shooting 
    private projectiles: Projectile[] = [];
    private lastShotTime: number = 0;
    private shootCooldown: number = 10;// was 250

    //variables for the enemies
    private enemies: EnemyType[] = [];
    private enemyVelocity: Velocity = {x:2, y:0};
    private enemyDropDistance: number = 40;

    private particles: Particle[] = [];
    private game = {
        over: false,
        active: false
    };

    //variables for the game play loop
    //this needs to be a hashmap idk how to do that lmao
    private levelSpawntimes = new Map<number, number>([[0,0],[1,5000]]); // come back to this how to properly type cast this 

    private currentLevel: number = 0;
    private isSpawning: boolean = false;


    public static playerImage: HTMLImageElement;
    public static enemyImage: HTMLImageElement;
    public static hunterImage: HTMLImageElement;
    private imagesLoaded: Promise<void>;


    constructor(canvas: HTMLCanvasElement, onMessage: (message: string) => void, onScoreChange?: (score: number) => void, onGameOver?: () => void) {
        this.canvas = canvas;
        const context = canvas.getContext("2d");
        this.onScoreChange = onScoreChange;
        this.onGameOver = onGameOver;
        this.onMessage = onMessage;
        this.score = 0;

        if (!context) throw new Error("could not get context");
        this.c = context;

        this.canvas.width = 1024;
        this.canvas.height = 576;

        this.imagesLoaded = this.loadAllImages().then(() => {
            this.player = new Player(this.c, this.canvas.width, this.canvas.height);

        });
        
        this.setupEventListeners();
    }

    private async loadAllImages(): Promise<void> {

        SpaceInvaders.playerImage = new Image();
        SpaceInvaders.enemyImage = new Image();
        SpaceInvaders.hunterImage = new Image();


        const playerPromise = new Promise<void>((resolve, reject) => {
            SpaceInvaders.playerImage.onload = () => {
                console.log('player image loaded');
                resolve();
            };
            SpaceInvaders.playerImage.onerror = () => reject(new Error('failed to load player'));
            SpaceInvaders.playerImage.src = playerImageUrl;
        });

        const enemyPromise = new Promise<void>((resolve, reject) => {
            SpaceInvaders.enemyImage.onload = () => {
                console.log('enemy image loaded');
                resolve();
            };
            SpaceInvaders.enemyImage.onerror = () => reject(new Error('failed to load enemy'));
            SpaceInvaders.enemyImage.src = enemyImg;


        });

        const hunterPromise = new Promise<void>((resolve, reject) => {
            SpaceInvaders.hunterImage.onload = () => {
                console.log('hunter enemy loaded');
                resolve();
            };
            SpaceInvaders.hunterImage.onerror = () => reject(new Error('failed to load hunter enemy'));
            SpaceInvaders.hunterImage.src = hunterImg;

        });

        await Promise.all([playerPromise, enemyPromise, hunterPromise]);
        console.log('images loaded success!')
    }

    public async waitForImages(): Promise<void> {
        await this.imagesLoaded;
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
            this.player.velocity.y = 0;
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
            // full controls I need to add a full mouse following as well

            if (this.score >= 3900) {
                if (this.keys['ArrowUp'] || this.keys['w']) {
                    this.player.velocity.y = -speed;
                    //this.player.rotation = -.2;
                }
                if (this.keys['ArrowDown'] || this.keys['s']) {
                    this.player.velocity.y = speed;
                    //this.player.rotation = .2;
                }
                if(this.keys[' ']) { // change this to take in click 
                    const currentTime = Date.now();
                    if (currentTime - this.lastShotTime >= this.shootCooldown){
                        this.lastShotTime = currentTime;
                        this.projectiles.push(new Projectile(this.c, {x: this.player.position.x + this.player.width / 2 ,y: this.player.position.y + 30 }, {x: 0 ,y: -1  }));


                    }
                }
            }
        }   

    }

    private spawnEnemyGrid(): void {
        const rows: number = 3;
        const cols: number = 13;
        const enemySpacing: number = 30;
        const startX: number = 100;
        const startY: number = 50;

        for (let row = 0; row < rows; row++){
            for (let col = 0; col < cols; col++) {
                this.enemies.push(new Enemy(this.c, {x: startX + col * enemySpacing, y: startY + row * enemySpacing}));
            }
        }
        //testing
        //for (let i = 0; i < 5; i++){
          //  setTimeout(() => this.onMessage('INCOMING HOSTILES DETECTED'), 1500);
        //}

    }
    private spawnHunters(): void {
        const min: number = 1;
        const max = 3;
        const spawnBoxWidth = 1024; //this is the canvas width as well
        let total = 0;

        total = Math.floor(Math.random() * (max - min + 1) + min);
        for (let count = 0; count < total; count++){
            let xPostiion = Math.random() * spawnBoxWidth;
            this.enemies.push(new HunterEnemy(this.c, {x: xPostiion, y: 0}, this.player))
            console.log('hunter spaned at:' , xPostiion)
        }


    }

    private updateEnemies(): void {
        let hitEdge = false;

        this.enemies.forEach(enemy => {
            if (enemy instanceof Enemy) {
                if ((enemy.position.x + enemy.width >= this.canvas.width && this.enemyVelocity.x > 0) || (enemy.position.x <= 0 && this.enemyVelocity.x < 0)) {
                    hitEdge = true;
                }
            }
        });

        if (hitEdge) {
            this.enemyVelocity.x *= -1;
            this.enemies.forEach(enemy => {
                if(enemy instanceof Enemy) {
                    enemy.position.y += this.enemyDropDistance;
                }
            });
        }
        this.enemies.forEach((enemy, i) => {
            if (enemy instanceof Enemy) {
                enemy.velocity = {...this.enemyVelocity};
            }
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

    createParticles(object: EnemyType | Player , color: string = 'purple'){
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

        //enemy spawing tying to make a game play loop and its fucking me 
        if (this.score >= 3800 && this.isSpawning == false) {
            this.isSpawning = true;
            if (this.currentLevel == 0){
                this.onMessage("well maybe not the classic, try WASD and moving the mouse");

                setTimeout(() => this.onMessage('Alright good luck!'), 3000);
            }
            setTimeout(() => this.spawnEnemies(), 6000);
            //setTimeout(() => this.spawnHunters(), spawnInterval);
            this.currentLevel = 1;

            setTimeout(() => this.isSpawning = false, this.levelSpawntimes.get(this.currentLevel));// come back this need to be a hash map
        }


    }
    start(): void {
        this.game.over = false;
        //testing 
        this.animate();
        this.spawnEnemyGrid();
        //this.spawnHunters();
        this.onMessage("Welcome to the classic Space Invaders");


        
        //this.gamePlayLoop();
                    
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

            if (this.onGameOver) {
                this.onGameOver();
            }
        }, 2000);
    }
    
    restart(): void {
        this.score = 0;
        this.game.over = false;
        this.game.active = true;

        this.projectiles = [];
        this.enemies = [];
        this.particles = [];

        this.player = new Player(this.c, this.canvas.width, this.canvas.height);

        this.keys = {};

        this.lastShotTime = 0;

        this.enemyVelocity = {x: 2, y: 0};

        if (this.onScoreChange) {
            this.onScoreChange(this.score);
        }
        this.setupEventListeners();
        this.start();

    }

    private spawnEnemies(): void {
        this.spawnEnemyGrid();
        this.spawnHunters();

        //this is being dealt with in the animate section now
        /*if (this.score >= 3900) {
            this.onMessage("well maybe not the classic, try WASD and moving the mouse");

            setTimeout(() => this.onMessage('Alright good luck!'), 3000);

            timeInterval = 5000;
        }
        return timeInterval; */


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

        this.image = SpaceInvaders.playerImage;

        this.width = this.image.width * 0.1;
        this.height = this.image.height * 0.1;

        this.position = { x: canvasWidth / 2 - this.width / 2, y: canvasHeight - this.height - 20 };
        this.velocity = { x: 0, y: 0 };
    }

    draw(): void {
        this.c.save()
        this.c.translate(this.position.x + this.width / 2, this.position.y + this.height / 2);
        this.c.rotate(this.rotation)
        this.c.translate(-this.position.x - this.width / 2, -this.position.y - this.height / 2);


        if (this.image && this.image.complete) {
            this.c.drawImage(this.image, this.position.x, this.position.y, this.width, this.height);
        }
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

        this.image = SpaceInvaders.enemyImage;

        this.width = this.image.width * 0.05;
        this.height = this.image.height * 0.05;
        this.velocity = { x: 0, y: 0 };
        this.position = startPosition || { x: 300, y: 300};

    }
    draw(): void {
        this.c.save();
        this.c.filter = 'hue-rotate(300deg) saturate(2) brightness(1.1)';
        if (this.image && this.image.complete) {
            this.c.drawImage(this.image, this.position.x, this.position.y, this.width, this.height);
        }
        this.c.restore();
    }
    update(): void {
        this.draw()
        this.position.x += this.velocity.x;
        this.position.y += this.velocity.y;
    }
}
class HunterEnemy {
    private c: CanvasRenderingContext2D;
    public position: Position;
    public velocity: Velocity;
    private image: HTMLImageElement;
    public readonly width: number;
    public readonly height: number;

    //state properties
    private state: 'hovering' | 'attacking' = 'hovering';
    private stateTimer: number = 0;
    private hoverDuration: number = 2000;
    private attackDuration: number = 500;
    private playerRef?: Player;
    private hoverSpeed: number = 1;
    private attackSpeed: number = 3;
    private rotation: number = 0;

    constructor(context: CanvasRenderingContext2D, startPosition?: Position, player?: Player) {
        this.c = context;
        this.image = SpaceInvaders.hunterImage;
        this.playerRef = player;
        this.stateTimer = Date.now();
        
        this.width = this.image.width * 0.05;
        this.height = this.image.height * 0.05;
        this.velocity = { x: 0, y: 0 };
        this.position = startPosition || { x: 200, y: 200};
    }

    draw(): void {
        this.c.save();
        this.c.translate( 
            this.position.x + this.width / 2,
            this.position.y + this.height / 2
        );
        if (this.state === 'attacking') {
            this.c.rotate(this.rotation);
        }

        this.c.translate(
            -this.position.x - this.width / 2, 
            -this.position.y - this.height / 2
        );

        if (this.image && this.image.complete) {
            this.c.drawImage(this.image, this.position.x, this.position.y, this.width, this.height);
        }
        this.c.restore();
    }
    update(): void {
        const currentTime = Date.now();
        const elapsedTime = currentTime - this.stateTimer;

        if (this.state === 'hovering') {
            this.velocity.x = Math.random() * Math.sin(currentTime * 0.002) * this.hoverSpeed;
            this.velocity.y = Math.random() * Math.cos(currentTime * 0.003) * this.hoverSpeed;

            if (elapsedTime >= this.hoverDuration) {
                this.state = 'attacking';
                this.stateTimer = currentTime;

                if (this.playerRef) {
                    const dx = this.playerRef.position.x - this.position.x;
                    const dy = this.playerRef.position.y - this.position.y;
                    const distance = Math.sqrt(dx * dx + dy * dy);
                    
                    this.velocity.x = (dx / distance) * this.attackSpeed;
                    this.velocity.y = (dy / distance) * this.attackSpeed;
                }
            }
        } else if (this.state === "attacking") {
            if (elapsedTime >= this.attackDuration) {
                this.state = 'hovering';
                this.stateTimer = currentTime;
            }
        }

        this.rotation = Math.atan2(this.velocity.y, this.velocity.x) - Math.PI / 2;

        this.draw();
        this.position.x += this.velocity.x;
        this.position.y += this.velocity.y;
    }
}