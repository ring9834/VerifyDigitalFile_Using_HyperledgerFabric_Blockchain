$(function(){
    // nav收缩展开
    $('.lnav-item>a').on('click',function(){
        if (!$('.lnav').hasClass('lnav-mini')) {
            if ($(this).next().css('display') == "none") {
                //展开未展开
                $('.lnav-item').children('ul').slideUp(300);
                $(this).next('ul').slideDown(300);
                $(this).parent('li').addClass('lnav-show').siblings('li').removeClass('lnav-show');
            }else{
                //收缩已展开
                $(this).next('ul').slideUp(300);
                $('.lnav-item.lnav-show').removeClass('lnav-show');
            }
        }
    });
    //nav-mini切换
    $('#mini').on('click',function(){
        if (!$('.lnav').hasClass('lnav-mini')) {
            $('.lnav-item.lnav-show').removeClass('lnav-show');
            $('.lnav-item').children('ul').removeAttr('style');
            $('.lnav').addClass('lnav-mini');
            $('#contentPage').css('width','calc(100% - 60px)');
            $('#contentPage').css('margin-left','60px');
        }else{
            $('.lnav').removeClass('lnav-mini');
            $('#contentPage').css('width','calc(100% - 220px)');
            $('#contentPage').css('margin-left','220px');
        }
    });
});