$(function(){
  $('#rssify').click(function(e){
    var input = $('#input');

    if (!input.val().match(/vk.com\/(\w+)/)) return;

    window.location.href = "http://rssify.me/api/v1/rss?g=" + input.val();

    e.preventDefault();
  })
});
