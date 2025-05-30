const burger = document.getElementById("burger");
const modal = document.getElementById("modal");
const clickBurger = document.getElementById("clickBurger");
const body = document.getElementsByName("body");

clickBurger.addEventListener("click", ()=> {
    modal.classList.toggle("modal-active");
    burger.classList.toggle("burger-active");
    $("body").toggleClass("no-scroll");
    $("html").toggleClass("no-scroll");
})

const modalItems = document.querySelectorAll('.modalItem');

modalItems.forEach(item => {
  item.addEventListener('click', () => {
    modal.classList.toggle("modal-active");
    burger.classList.toggle("burger-active");
    $("body").toggleClass("no-scroll");
    $("html").toggleClass("no-scroll");
  });
});

let currentSlide = 0;
const slides = document.querySelectorAll('.slide');

function showSlide(index) {
    slides.forEach((slide, i) => {
        if (i === index) {
            slide.classList.add('slide-active');
        } else {
            slide.classList.remove('slide-active');
        }
    });
}

function changeSlide(direction) {
    currentSlide += direction;
    if (currentSlide < 0) {
        currentSlide = slides.length - 1;
    } else if (currentSlide >= slides.length) {
        currentSlide = 0;
    }
    showSlide(currentSlide);
}
showSlide(currentSlide);

// slider1
$(document).ready(function(){
    $('.slider1').slick({
      slidesToShow: 5, // 5 полных слайда
      slidesToScroll: 1, // Перелистывание по 1 слайду
      centerMode: true, // Центрирование слайдов
      variableWidth: true, // Частично видимые слайды
      autoplay: true, // Автопрокрутка
      autoplaySpeed: 2000, // Интервал 5 секунд
      arrows: false, // Убираем стрелки
      responsive: [
        {
          breakpoint: 1024,
          settings: {
            slidesToShow: 3,
          }
        },
        {
          breakpoint: 768,
          settings: {
            slidesToShow: 2,
          }
        },
        {
          breakpoint: 480,
          settings: {
            slidesToShow: 1,
          }
        }
      ]
    });
  });


// slider2
$(document).ready(function(){
    $('.slider2').slick({
      slidesToShow: 5, // 5 полных слайда
      slidesToScroll: 1, // Перелистывание по 1 слайду
      centerMode: true, // Центрирование слайдов
      variableWidth: true, // Частично видимые слайды
      autoplay: true, // Автопрокрутка
      autoplaySpeed: 800, // Интервал 5 секунд
      arrows: false, // Убираем стрелки
      responsive: [
        {
          breakpoint: 1024,
          settings: {
            slidesToShow: 3,
          }
        },
        {
          breakpoint: 768,
          settings: {
            slidesToShow: 2,
          }
        },
        {
          breakpoint: 480,
          settings: {
            slidesToShow: 1,
          }
        }
      ]
    });
  });



// --Form--

const formPopup = document.getElementById('feedback-form');
const form = document.getElementById('feedback-form-data');
const feedbackMessage = document.getElementById('feedback-message');
const close = document.getElementById('close');


function saveFormData() {
  localStorage.setItem('formData', JSON.stringify({
    name: document.getElementById('name').value,
    email: document.getElementById('email').value,
    phone: document.getElementById('phone').value,
    company: document.getElementById('company').value,
    message: document.getElementById('message').value,
  }));
}

function restoreFormData() {
  const formData = localStorage.getItem('formData');
  if (formData) {
    const data = JSON.parse(formData);
    document.getElementById('name').value = data.name;
    document.getElementById('email').value = data.email;
    document.getElementById('phone').value = data.phone;
    document.getElementById('company').value = data.company;
    document.getElementById('message').value = data.message;
  }
}


// --FAQ
const faqItems = document.querySelectorAll('.faq-item');

faqItems.forEach(item => {
  item.addEventListener('click', () => {
    item.classList.toggle('active');
  });
});
