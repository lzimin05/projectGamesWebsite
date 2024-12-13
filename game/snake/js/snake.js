const canvas = document.getElementById("game");
const ctx = canvas.getContext("2d");

const ground = new Image();
ground.src = "game/snake/image/ground.png";

const foodImg = new Image();
foodImg.src = "game/snake/image/apple3.png";

let box = 32;
let speed_game = 150;
let score = 0;

let food = {
    x: Math.floor((Math.random() * 17 + 1)) * box,
    y: Math.floor((Math.random() * 15 + 3)) * box,
};

let snake = [];
snake[0] = {
    x: 9 * box,
    y: 10 * box
};

let retry = document.querySelector('.game_try_again');
let dir;

function direction(event) {
    if (((event.keyCode == 37) || (event.keyCode == 65)) && dir != "right")
        dir = "left";
    else if (((event.keyCode == 38) || (event.keyCode == 87)) && dir != "down")
        dir = "up";
    else if (((event.keyCode == 39) || (event.keyCode == 68)) && dir != "left")
        dir = "right";
    else if (((event.keyCode == 40) || (event.keyCode == 83)) && dir != "up")
        dir = "down";
}

function eatTail(head, arr) {
    for (let i = 0; i < arr.length; i++) {
        if (head.x == arr[i].x && head.y == arr[i].y){
            clearInterval(game);
            button_retry.classList.remove('button_hide');
            label_score.textContent=`GameOver \n Your Score: ${score}`;
        }
    }
}

function drawGame() {
    console.log("1")
    label_score.textContent=''
    ctx.drawImage(ground, 0, 0);

    ctx.drawImage(foodImg, food.x, food.y);

    for (let i = 0; i < snake.length; i++) {
        ctx.fillStyle = i == 0 ? "green" : "red";
        ctx.fillRect(snake[i].x, snake[i].y, box, box);
    }

    ctx.fillStyle = "white";
    ctx.font = "50px Arial";
    ctx.fillText(score, box * 2.5, box * 1.7);

    let snakeX = snake[0].x;
    let snakeY = snake[0].y;

    if (snakeX == food.x && snakeY == food.y) {
        generateFood(); //генерация еды
    } else {
        snake.pop();
    }

    if (snakeX < box || snakeX > box * 17 ||	//GameOver
        snakeY < 3 * box || snakeY > box * 17) {
		//alert(`GameOver\nYour score: ${score}`)
        label_score.textContent=`GameOver \n Your Score: ${score}`;
        clearInterval(game)
		button_retry.classList.remove('button_hide')
	}
    if (dir == "left") snakeX -= box;
    if (dir == "right") snakeX += box;
    if (dir == "up") snakeY -= box;
    if (dir == "down") snakeY += box;

    let newHead = {
        x: snakeX,
        y: snakeY
    };

    eatTail(newHead, snake);

    snake.unshift(newHead);
}

function generateFood() {
    score++;
    if (score == 255) { //если победа
        label_score.textContent=`YOU WIN!!! \n Your score: ${score}`
        clearInterval(game);
        button_retry.classList.remove('button_hide');

    }
    food = {
        x: Math.floor((Math.random() * 17 + 1)) * box, //17x15 = 255
        y: Math.floor((Math.random() * 15 + 3)) * box,
    };
    while (true) {
        for(var i = 0; i < snake.length; i++) {
            if ((snake[i].x == food.x) && (snake[i].y == food.y)) {
                if(food.x < 18*32) { //17+1 (Xmax)
                    food.x = food.x + box;
                } else if(y < 18*32) { //15+3 (Ymax)
                    food.y = food.y + box;
                    food.x = 1 * box;
                } else {
                    food.y = food.y + 3 * box;
                    food.x = food.x + 1 * box;
                }
                i = 0;
                continue;
            }
        }
        break;
    }
}

function newGame() {
    dir="up"
    label_score.textContent='';
    console.log("2");
    snake.length = 0;
    head = {
        x: 9 * box,
        y: 10 * box
    };
    snake.push(head);
    button_retry.classList.add('button_hide');
    button_retry.disabled = true;
    setTimeout(() => button_retry.disabled = false, 3000);
    setTimeout(() => {game = setInterval(drawGame, speed_game);}, 1000); //предотвратить двойного нажатия на кнопку
    score = 0;

}   

let button_retry = document.querySelector('.retry');

if (button_retry.classList.contains('button_hide')) {
    document.addEventListener("keydown", direction);
}

console.log(button_retry)
button_retry.addEventListener("click", newGame);
label_score = document.querySelector('.score');

let game = setInterval(drawGame, speed_game); //скорость игры
