
// make the icon glow when the input field is focused
$(".input_text").focus(function () { // focus is when the mouse is on the input field
    $(this).prev('.fa').addClass('glowIcon');
});

$(".input_text").focusout(function () {
    $(this).prev('.fa').removeClass('glowIcon');
});

/* --------------------------------------- Login ------------------------------------------ */

//change the submit button color when the mouse is over it
$("#login_button").hover(function () {
    $(this).addClass('button_group');
});

$("#login_button").mouseleave(function () {
    setTimeout(() => {
        $(this).removeClass('button_group');
    }, 100);
});

/* --------------------------------------- Register ------------------------------------------ */

$("#register_button").hover(function () {
    $(this).addClass('button_group');
});

$("#register_button").mouseleave(function () {
    setTimeout(() => {
        $(this).removeClass('button_group');
    }, 100);
});

/* ----------------------------------- Age style ---------------------------------------------- */

jQuery('<div class="quantity-nav"><div class="quantity-button quantity-up">+</div><div class="quantity-button quantity-down">-</div></div>').insertAfter('.quantity input');
jQuery('.quantity').each(function () {
    var spinner = jQuery(this),
        input = spinner.find('input[type="number"]'),
        btnUp = spinner.find('.quantity-up'),
        btnDown = spinner.find('.quantity-down'),
        min = input.attr('min'),
        max = input.attr('max'),
        age = input.attr('placeholder');

    btnUp.click(function () {
        var oldValue = parseFloat(input.val());
        if (isNaN(oldValue)) {
            var newVal = 0;
        } else if (oldValue >= max) {
            var newVal = oldValue;
        } else {
            var newVal = oldValue + 1;
        }
        spinner.find("input").val(newVal);
        spinner.find("input").trigger("change");
    });

    btnDown.click(function () {
        var oldValue = parseFloat(input.val());
        if (isNaN(oldValue)) { // if the age isn't changed, set the value to 0
            var newVal = 0;
        } else if (oldValue <= min) {
            var newVal = oldValue;
        } else {
            var newVal = oldValue - 1;
        }
        spinner.find("input").val(newVal);
        spinner.find("input").trigger("change");
    });
});

/* ----------------------------------- Forms Fade ---------------------------------------------- */

// fade in the register form and fade out the login form

$("#Sign_Up").click(function () {
    $(".login_form_container").fadeOut(fade_out, function () {
        $(".register_form_container").fadeIn(fade_in);
    });
});

// fade in the login form and fade out the register form

$("#Sign_In").click(function () {
    $(".register_form_container").fadeOut(fade_out, function () {
        $(".login_form_container").fadeIn(fade_in);
    });
});

// get the fade_out keyframe from the css file
var fade_out = getComputedStyle(document.documentElement).getPropertyValue('--fade-out');
var fade_in = getComputedStyle(document.documentElement).getPropertyValue('--fade-in');

/* ------------------------------------- Forum Header ----------------------------------------------------- */

//change the logout button color when the mouse is over it
$("#logout_button").hover(function () {
    $(this).addClass('button_group');
});

$("#logout_button").mouseleave(function () {
    setTimeout(() => {
        $(this).removeClass('button_group');
    }, 100);
});



/* -----------------------------------Connected User---------------------------------------------- */

function enterForum () {
    $(".login_form_container").fadeOut(fade_out, function () {
    
        $(".global_container").fadeIn(fade_in);
        $(".ppl_container").fadeIn(fade_in);
        $(".forum_container").fadeIn(fade_in);
        
    });
}

/* ----------------------------------- Chat ---------------------------------------------- */

// open the chat box
$("#chat_button").click(function () {
    $(".chat_container").fadeIn(fade_in);
});


/* -----------------------------------Private messages---------------------------------------------- */

/* -----------------------------------Create a Post---------------------------------------------- */
//button
$("#create_post").hover(function () {
    $(this).addClass('button_group');
});

$("#create_post").mouseleave(function () {
    setTimeout(() => {
        $(this).removeClass('button_group');
    }, 100);
});

// open the create post box
$("#create_post").click(function () {
    $("#new_post").fadeIn(fade_in);
});




