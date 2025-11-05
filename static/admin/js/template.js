(function($) {
  'use strict';
  $(function() {
    var body = $('body');
    var contentWrapper = $('.content-wrapper');
    var scroller = $('.container-scroller');
    var footer = $('.footer');
    var sidebar = $('.sidebar');

    //Add active class to nav-link based on url dynamically
    function addActiveClass(element) {
      if (current === "") {
        if (element.attr('href').indexOf("index.html") !== -1) {
          element.parents('.nav-item').last().addClass('active');
          if (element.parents('.sub-menu').length) {
            element.closest('.collapse').addClass('show');
            element.addClass('active');
          }
        }
      } else {
        if (element.attr('href').indexOf(current) !== -1) {
          element.parents('.nav-item').last().addClass('active');
          if (element.parents('.sub-menu').length) {
            element.closest('.collapse').addClass('show');
            element.addClass('active');
          }
          if (element.parents('.submenu-item').length) {
            element.addClass('active');
          }
        }
      }
    }

    var current = location.pathname.split("/").slice(-1)[0].replace(/^\/|\/$/g, '');
    $('.nav li a', sidebar).each(function() {
      var $this = $(this);
      addActiveClass($this);
    });

    //Close other submenu in sidebar on opening any
    sidebar.on('show.bs.collapse', '.collapse', function() {
      sidebar.find('.collapse.show').collapse('hide');
    });

    //Change sidebar 
    $('[data-toggle="minimize"]').on("click", function() {
      body.toggleClass('sidebar-icon-only');
    });

    //checkbox and radios
    $(".form-check label,.form-radio label").append('<i class="input-helper"></i>');
  });

  // ✅ تأكيد وجود العناصر قبل التعديل
  const proBanner = document.querySelector('#proBanner');
  const navbar = document.querySelector('.navbar');
  const pageBodyWrapper = document.querySelector('.page-body-wrapper');
  const bannerClose = document.querySelector('#bannerClose');

  $('#navbar-search-icon').click(function() {
    $("#navbar-search-input").focus();
  });

  if ($.cookie('royal-free-banner') != "true") {
    if (proBanner) proBanner.classList.add('d-flex');
    if (navbar) navbar.classList.remove('fixed-top');
  } else {
    if (proBanner) proBanner.classList.add('d-none');
    if (navbar) navbar.classList.add('fixed-top');
  }

  if (navbar && pageBodyWrapper) {
    if ($(navbar).hasClass("fixed-top")) {
      pageBodyWrapper.classList.remove('pt-0');
      navbar.classList.remove('pt-5');
    } else {
      pageBodyWrapper.classList.add('pt-0');
      navbar.classList.add('pt-5');
      navbar.classList.add('mt-3');
    }
  }

  if (bannerClose) {
    bannerClose.addEventListener('click', function() {
      if (proBanner) {
        proBanner.classList.add('d-none');
        proBanner.classList.remove('d-flex');
      }
      if (navbar && pageBodyWrapper) {
        navbar.classList.remove('pt-5');
        navbar.classList.add('fixed-top');
        pageBodyWrapper.classList.add('proBanner-padding-top');
        navbar.classList.remove('mt-3');
      }
      var date = new Date();
      date.setTime(date.getTime() + 24 * 60 * 60 * 1000);
      $.cookie('royal-free-banner', "true", { expires: date });
    });
  }

})(jQuery);
