import { Play } from "lucide-svelte";
import playerImageUrl from './imgs/player.png';
import { siVelocity } from "simple-icons";
import enemyImg from './imgs/enemy.png';
import hunterImg from './imgs/enemyHunter.png'
import rangedImg from './imgs/rangedEnemy.png'

interface Position {
    x: number;
    y: number;
}

interface Velocity {
    x: number;
    y: number;
}
type EnemyType = Enemy | HunterEnemy | RangedEnemy

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
    private mouseDown: boolean = false;
    private mousePosition: Position = { x:0, y:0};

    //variables for the shooting 
    private projectiles: Projectile[] = [];
    private lastShotTime: number = 0;
    private shootCooldown: number = 250;// was 250

    //variables for the enemies
    private enemies: EnemyType[] = [];
    private enemyVelocity: Velocity = {x:2, y:0};
    private enemyDropDistance: number = 40;
    private spawnRate: number = 1; // this is what is controling how hard the game is I believe

    private enemyProjectiles: Projectile[] = [];
    private particles: Particle[] = [];
    private game = {
        over: false,
        active: false
    };

    //variables for the game play loop
    //this needs to be a hashmap idk how to do that lmao
    private levelSpawntimes = new Map<number, number>([[0,0],[1,10000],[2,8000]]); // come back to this how to properly type cast this 

    private currentLevel: number = 0;
    private isSpawning: boolean = false;
    private spawnTimeoutId: any;


    public static playerImage: HTMLImageElement;
    public static enemyImage: HTMLImageElement;
    public static hunterImage: HTMLImageElement;
    public static rangedImage: HTMLImageElement;
    private imagesLoaded: Promise<void>;

    // debugging frames per second type shit
    private fps = 0;

    //for throttling the game at 60 fps
    private targetFPS = 100;
    private frameInterval = 1000 / this.targetFPS;
    private lastFrameTime = 0;


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
        SpaceInvaders.rangedImage = new Image();


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
        const rangedPromise = new Promise<void>((resolve, reject) => {
            SpaceInvaders.rangedImage.onload = () => {
                resolve();
            };
            SpaceInvaders.rangedImage.onerror = () => reject(new Error('failed to load ranged enemy'));
            SpaceInvaders.rangedImage.src = rangedImg;
        })

        await Promise.all([playerPromise, enemyPromise, hunterPromise, rangedPromise]);
        console.log('images loaded success!')
    }

    public async waitForImages(): Promise<void> {
        await this.imagesLoaded;
    }

    private setupEventListeners(): void {
        window.addEventListener('keydown', this.handleKeyDown);
        window.addEventListener('keyup', this.handleKeyUp);
        this.canvas.addEventListener('mousedown', this.handleMouseDown);
        this.canvas.addEventListener('mouseup', this.handleMouseUp);
        this.canvas.addEventListener('mousemove', this.handleMouseMove);
    }
    private handleKeyDown = (event: KeyboardEvent): void => {
        this.keys[event.key] = true;
    }
    private handleKeyUp = (event: KeyboardEvent): void => {
        this.keys[event.key] = false;
    }
    private handleMouseDown = (): void => {
        this.mouseDown = true;
    }
    private handleMouseUp = (): void => {
        this.mouseDown = false;
    }
    private handleMouseMove = (event: MouseEvent): void => {
        const rect = this.canvas.getBoundingClientRect();
        this.mousePosition.x = event.clientX - rect.left;
        this.mousePosition.y = event.clientY - rect.top;
    }

    private updatePlayerControls(): void {
        if (!(this.game.over)) {
            this.player.rotation = 0;
            this.player.velocity.x = 0;
            this.player.velocity.y = 0;
            const speed = 5;
            

            if (this.keys['ArrowLeft'] || this.keys['a']) {
                if (this.player.position.x > 0 - (this.player.width / 1.75)) {
                    this.player.velocity.x = -speed;
                    this.player.rotation = -.2;
                }
            }
            if (this.keys['ArrowRight'] || this.keys['d']) {
                if (this.player.position.x < this.canvas.width - (this.player.width / 1.75)) {
                    this.player.velocity.x = speed;
                    this.player.rotation = .2;
                }
            }
            if(this.keys[' '] || this.mouseDown && this.score < 3000) {
                const currentTime = Date.now();
                if (currentTime - this.lastShotTime >= this.shootCooldown){
                    this.lastShotTime = currentTime;
                    this.projectiles.push(new Projectile(this.c, {x: this.player.position.x + this.player.width / 2 ,y: this.player.position.y + 30 }, {x: 0 ,y: -1  }));
                }
            }
            // full controls 

            if (this.score >= 3000) {
                const dx = this.mousePosition.x - (this.player.position.x + this.player.width / 2);
                const dy = this.mousePosition.y - (this.player.position.y + this.player.height / 2);
                const angleToMouse = Math.atan2(dy, dx);

                this.player.rotation = angleToMouse + Math.PI / 2; 
                if (this.keys['ArrowUp'] || this.keys['w']) {
                    if (this.player.position.y > 0 - (this.player.height / 1.75)) {
                        this.player.velocity.y = -speed;
                        //this.player.rotation = -.2;
                    }
                }
                if (this.keys['ArrowDown'] || this.keys['s']) {
                    if (this.player.position.y < this.canvas.height - (this.player.height / 1.75)) {
                        this.player.velocity.y = speed;
                        //this.player.rotation = .2;
                    }
                }
                if(this.keys[' '] || this.mouseDown) {
                    const currentTime = Date.now();
                    if (currentTime - this.lastShotTime >= this.shootCooldown){
                        this.lastShotTime = currentTime;
                        const shootAngle = this.player.rotation - Math.PI / 2;

                        //const shootVelocity = { x: Math.cos(shootAngle), y: Math.sin(shootAngle) };
                        this.projectiles.push(new Projectile(this.c, {x: this.player.position.x + this.player.width / 2 ,y: this.player.position.y + 30 }, { x: Math.cos(shootAngle), y: Math.sin(shootAngle) }));


                    }
                }
            }
        }   

    }

    private spawnEnemyGrid(spawnRate: number = 1): void {
        let rows: number = Math.round(3 * spawnRate);
        let cols: number = Math.round(10 * spawnRate);

        const spawnBoxWidth = 1024 //this is the canvas width

        const enemySpacing: number = 30;
        const startY: number = 20;

        const gridWidth = cols * enemySpacing;

        let maxStartX = spawnBoxWidth - gridWidth - 20;
        let minStartX = 20;

        let startX: number = Math.random() * (maxStartX - minStartX) + minStartX;

        for (let row = 0; row < rows; row++){
            for (let col = 0; col < cols; col++) {
                this.enemies.push(new Enemy(this.c, {x: startX + col * enemySpacing, y: startY + row * enemySpacing}));
            }
        }

    }
    private spawnHunters(spawnRate: number = 1): void {
        let min: number = Math.round((2 * spawnRate) + this.currentLevel);
        let max: number = Math.round((5 * spawnRate) + this.currentLevel);
        const spawnBoxWidth = 1024; //this is the canvas width as well
        let total = 0;

        total = Math.floor(Math.random() * (max - min + 1) + min);
        for (let count = 0; count < total; count++){
            let xPosition = Math.random() * spawnBoxWidth;
            this.enemies.push(new HunterEnemy(this.c, {x: xPosition, y: 0}, this.player))
            //console.log('hunter spaned at:' , xPosition)
        }
    }

    private spawnRangedEnemies(spawnRate: number = 1): void {
        let min: number = Math.round(1 * spawnRate);
        let max: number = Math.round(3 * spawnRate);
        let total = 0;
        total = Math.floor(Math.random() * (max - min + 1) + min);

        for (let count = 0; count < total; count++) {
            let randomX = Math.random() * this.canvas.width;
            let randomY = Math.random() * this.canvas.height;

            const chosenSide = Math.random() < 0.5 ? randomX : randomY;
            if (chosenSide == randomX) {
                let yPosition = Math.random() < 0.5 ? 0 : this.canvas.height;
                this.enemies.push(new RangedEnemy(this.c, {x: randomX, y: yPosition}, this.player))
            } else if (chosenSide == randomY) {
                let xPosition = Math.random() < 0.5 ? 0 : this.canvas.width;
                this.enemies.push(new RangedEnemy(this.c, {x: xPosition, y: randomY}, this.player ))
            }
            console.log('ranged enemy spawned')

        }
    }

    private updateEnemies(): void {
        let hitEdge = false;

        this.enemies.forEach((enemy, i) => {
            if (enemy instanceof Enemy && enemy.position.y > this.canvas.height) {
                this.enemies.splice(i, 1);
                //var enemies_test = this.enemies.length
                console.log("enemy deleted total enemies left: %d", this.enemies.length)
            }
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
            if (enemy instanceof RangedEnemy && enemy.shooting) {
                const shootAngle = Math.atan2(
                    this.player.position.y - enemy.position.y,
                    this.player.position.x - enemy.position.x
                );
                this.enemyProjectiles.push(new Projectile(
                    this.c,
                    { x: enemy.position.x + enemy.width / 2, y: enemy.position.y + enemy.height / 2},
                    { x: Math.cos(shootAngle), y: Math.sin(shootAngle) },
                    'blue'
                ));
                this.createParticles(enemy, 'blue')
                enemy.shooting = false;
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
    
    animate = (currentTime: number = 0): void => {


        this.animationId = requestAnimationFrame(this.animate);
        const deltaTime = currentTime - this.lastFrameTime;
        if (deltaTime < this.frameInterval) {
            return
        }
        this.lastFrameTime = currentTime - (deltaTime % this.frameInterval)


        this.c.fillStyle = 'black';
        this.c.fillRect(0,0, this.canvas.width, this.canvas.height);
        this.updatePlayerControls();

        // FPS is stuff 
        this.fps = Math.round(1000 / deltaTime)
        this.c.fillStyle = 'white';
        this.c.font = '20px Arial';
        this.c.fillText(`FPS: ${this.fps}`, 10, 30);

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
            if (projectile.position.y > this.canvas.height || 
                projectile.position.x < 0 || 
                projectile.position.x > this.canvas.width) {
                setTimeout(() => {
                this.enemyProjectiles.splice(idx, 1);
            }, 0);
            } else {
                projectile.update();
            }
        });  
        this.enemyProjectiles.forEach((projectile, idx) => {
            if (projectile.position.y > this.canvas.height || 
                projectile.position.x < 0 || 
                projectile.position.x > this.canvas.width) {
                setTimeout(() => {
                this.enemyProjectiles.splice(idx, 1);
            }, 0);
            } else {
                projectile.update();
            }

            if (
                projectile.position.y - projectile.radius <= this.player.position.y + this.player.height &&
                projectile.position.x + projectile.radius >= this.player.position.x &&
                projectile.position.x - projectile.radius <= this.player.position.x + this.player.width &&
                projectile.position.y + projectile.radius >= this.player.position.y
            ) {
                console.log("Player hit!");
                this.player.isAlive = false;
                if (!this.game.over) {
                    for (let i = 0; i < 15; i++) {
                        this.createParticles(this.player, 'grey');
                    }
                }
                this.stop();
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
        if (this.score >= 3000 && !this.isSpawning) {
            this.isSpawning = true;
            if (this.currentLevel == 0){
                this.onMessage("well maybe not the classic, try WASD and moving the mouse");

                setTimeout(() => this.onMessage('Alright good luck!'), 3000);
                this.currentLevel = 1;
                setTimeout(() => this.onMessage('See how far you can get! Check the leaderboards below to see the competition'), 6000);

            } else if (this.currentLevel == 1 && this.score >= 8000) {
                this.onMessage("Nice you are getting the hang of it");
                setTimeout(() => this.onMessage("See you at the top of the leaderboard"), 3000);
                this.currentLevel = 2;
            }
            this.spawnLoop();
        }
    }
    spawnLoop() {
        this.spawnEnemies(this.spawnRate)

        this.spawnTimeoutId = setTimeout(() =>{
            this.spawnLoop();
            console.log("spawing out these johns")
        }, this.levelSpawntimes.get(this.currentLevel));

    }
    start(): void {
        this.game.over = false;
        //testing 
        this.animate();
        this.spawnEnemyGrid();
        //this.spawnRangedEnemies();
        //this.spawnHunters();
        this.onMessage("Welcome to the classic Space Invaders");                    
    }
    stop(): void {
        this.game.over = true;


        setTimeout(() => {
            this.spawnRate = 1;
            this.currentLevel = 0;
            clearTimeout(this.spawnTimeoutId);
            this.isSpawning = false;

        
            if (this.animationId != null) {
                cancelAnimationFrame(this.animationId);
                this.animationId = null;
            }

            window.removeEventListener('keydown', this.handleKeyDown);
            window.removeEventListener('keyup', this.handleKeyUp);
            this.canvas.removeEventListener('mousedown', this.handleMouseDown);
            this.canvas.removeEventListener('mouseup', this.handleMouseUp);
            this.canvas.removeEventListener('mousemove', this.handleMouseMove); 


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
        this.enemyProjectiles = [];

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

    private spawnEnemies(spawnRate: number): void {
        this.spawnRate *= 1.03
        this.spawnEnemyGrid(spawnRate);
        this.spawnHunters(spawnRate);
        if (this.score >= 8000) {
            this.spawnRangedEnemies(spawnRate)
        }

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

        this.width = this.image.width * 1.5;
        this.height = this.image.height * 1.5;

        this.position = { x: canvasWidth / 2 - this.width / 2, y: canvasHeight - this.height - 20 };
        this.velocity = { x: 0, y: 0 };
    }

    draw(): void {
        this.c.save()
        this.c.translate(this.position.x + this.width / 2, this.position.y + this.height / 2);
        this.c.rotate(this.rotation)
        this.c.translate(-this.position.x - this.width / 2, -this.position.y - this.height / 2);


        if (this.image && this.image.complete) {
            this.c.imageSmoothingEnabled = false;
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
    private color: string;

    constructor(context: CanvasRenderingContext2D, position: Position, velocity: Velocity, color: string = 'red' ) {
        this.position = position;
        this.velocity = velocity;
        this.c = context;
        this.color = color
        
        this.speed = 5;
        this.radius = 3;

    }
    draw() {
        this.c.beginPath();
        this.c.arc(this.position.x, this.position.y, this.radius, 0, Math.PI * 2);
        this.c.fillStyle = this.color;
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

        this.width = this.image.width;
        this.height = this.image.height;
        this.velocity = { x: 0, y: 0 };
        this.position = startPosition || { x: 300, y: 300};

    }
    draw(): void {
        this.c.save();
        if (this.image && this.image.complete) {
            this.c.imageSmoothingEnabled = false;
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
        
        this.width = this.image.width * 1.5;
        this.height = this.image.height * 1.5;
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
            this.c.imageSmoothingEnabled = false;
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

class RangedEnemy {
    private c: CanvasRenderingContext2D;
    public position: Position;
    public velocity: Velocity;
    private image: HTMLImageElement;
    public readonly width: number;
    public readonly height: number;

    private state: 'hovering' | 'attacking' = 'hovering';
    private stateTimer: number = 0;
    private hoverDuration: number = 2000;
    private attackDuration: number = 500;
    private playerRef?: Player;
    private hoverSpeed: number = 1;
    //private attackSpeed: number = 3;
    private rotation: number = 0;
    public shooting: boolean = false;

    constructor(context: CanvasRenderingContext2D, startPosition: Position, player: Player) {
        this.c = context;
        this.image = SpaceInvaders.rangedImage;
        this.playerRef = player;
        this.width = this.image.width * 1.5;
        this.height = this.image.height * 1.5;
        this.velocity = { x: 0, y: 0 };
        this.position = startPosition || { x: 200, y: 200};
        this.stateTimer = Date.now();

    }

    draw(): void {
        this.c.save()
        this.c.translate (
            this.position.x + this.width / 2,
            this.position.y + this.height / 2
        );
        this.c.rotate(this.rotation)

        this.c.translate(
            -this.position.x - this.width / 2, 
            -this.position.y - this.height / 2
        );

        if (this.image && this.image.complete) {
            this.c.imageSmoothingEnabled = false;
            this.c.drawImage(this.image, this.position.x, this.position.y, this.width, this.height);
        }
        this.c.restore()

    }
    update(): void {
        const currentTime = Date.now()
        const elapsedTime = currentTime - this.stateTimer;
        if(this.playerRef){
            const dx = this.playerRef.position.x - this.position.x;
            const dy = this.playerRef.position.y - this.position.y;
            this.rotation = Math.atan2(dy, dx) - Math.PI / 2;
        }

        if (this.state === 'hovering') {
            this.velocity.x = Math.random() * Math.sin(currentTime * 0.001) * this.hoverSpeed;
            this.velocity.y = Math.random() * Math.cos(currentTime * 0.001) * this.hoverSpeed;

            if (elapsedTime >= this.hoverDuration){
                this.state = 'attacking';
                this.stateTimer = currentTime;
                this.shooting = true;
            }
        } else if (this.state === 'attacking') {
            this.velocity.x = 0;
            this.velocity.y = 0;
            if (elapsedTime >= this.attackDuration) {
                this.state = 'hovering';
                this.stateTimer = currentTime;
                this.shooting = false;
            }
        }
        this.position.x += this.velocity.x;
        this.position.y += this.velocity.y;
        if (this.position.x < 0) {
            this.position.x = 0;
        } else if (this.position.x > 1024) {
            this.position.x = 1024;
        }

        if (this.position.y < 0)  {
            this.position.y = 0;
        } else if (this.position.y > 576) {
            this.position.y = 576
        }

        this.draw()

    }

}