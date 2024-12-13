//header
let button_selected = document.querySelectorAll('.menu');   //указывает на какой кнопке сейчас фокус
button_selected[0].className += ' selected';
let id_button_selected = 0;
//button_selected[0].documentElement.style.setProperty('.selected', nil)
for(let i = 0; i < button_selected.length; i++) {
    button_selected[i].addEventListener("click", function() {                   
        let ClassList = button_selected[i].classList;
        if (!ClassList.contains('selected')) {
            button_selected[id_button_selected].classList.remove('selected');   
            button_selected[i].className += ' selected';
            id_button_selected = i;
        }   
    });
}
//Обработка кнопки смены темы (светлая или темная)
let button_theme = document.querySelector('.dark_mode_button');
document.body.className += 'light_mode';
button_theme.addEventListener("click", function() {
    let ClassList = document.body.classList;
    if (ClassList.contains('light_mode')) {
        document.body.classList.remove('light_mode');
        document.body.classList.add('dark_mode');
    } else {
        document.body.classList.add('light_mode');
        document.body.classList.remove('dark_mode');
    }
});

//спрятать кнопку search
function hide_search_button() {
    if ((search_button.classList.contains('button_hide')) && (input_label.value != '')) {
        //console.log("a")
        search_button.classList.remove('button_hide');
        cansel_button_for_search.classList.remove('button_hide');
    } else if ((input_label.value == '') && (!search_button.classList.contains('button_hide'))){
        //console.log("b")
        search_button.classList.add('button_hide')
        cansel_button_for_search.classList.add('button_hide');
    }
}

let search_button = document.querySelector('.center_menu_search_button');
let input_label = document.querySelector('.center_menu_search_label');
let cansel_button_for_search = document.querySelector('.center_menu_cansel')

input_label.addEventListener('click', hide_search_button)
input_label.addEventListener('blur', hide_search_button)
input_label.addEventListener('input', hide_search_button)

//если нажали на крестик
let button_delete = document.querySelector('.center_menu_cansel');

button_delete.addEventListener('click', function(){
    input_label.value = '';
    hide_search_button();
})

//поиск карточки
//search_button input_label

let cards = document.querySelectorAll('.game_card') 
console.log(cards)
let items = ["Змейка", "Сокобан", "Ну, погоди!", "трафик райдер 3D"] //все игры
console.log(items);
search_button.addEventListener('click', search_cards);
document.addEventListener( 'keyup', event => {
    if( event.code === 'Enter' ) {
        console.log('enter was pressed');
        search_cards();
    }
});
button_delete.addEventListener('click', search_cards)

function search_cards() {
    document.querySelector('.info_bar').classList.add('button_hide')
    document.querySelector('.info_bar').classList.add('zero_params')
    const query = input_label.value.toLowerCase();
    let cardsum = 0;
    for(let i = 0; i < items.length; i++) {
        if ((items[i].toLowerCase()).includes(query)) {
            cardsum++;
            cards[i].classList.remove('button_hide')
            cards[i].classList.remove('zero_params') 
        } else {
            cards[i].classList.add('button_hide')
            cards[i].classList.add('zero_params')  
        }
    }
    if (cardsum == 0) {
        document.querySelector('.info_bar').classList.remove('button_hide')
        document.querySelector('.info_bar').classList.remove('zero_params')
        for(let j = 0; j < cards.length; j++) {
            cards[j].classList.remove('button_hide')
            cards[j].classList.remove('zero_params') 
        }
    }
    if (query == '') {
        for(let i = 0; i < cards.length; i++) {
            cards[i].classList.remove('button_hide')
            cards[i].classList.remove('zero_params')
        }
    }
}
