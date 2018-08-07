(function () {
    "use strict";
    $(document).ready(function () {
        $(".file").change(function (e) {
            populateChannels(e);
        })
    });

    function populateChannels(e) {
        let file = $(e.target)
        var channels = $(".channel");
        var channel = null;
        if (file.hasClass("red")) {
            channel = channels.find(".red");
        }
        let fileSelection = file.val()
        let x = UploadData.Names;
        let selectedIdx = UploadData.Names.findIndex(function (elem) {
            return elem == file.val();
        })
        if (selectedIdx == -1) {
            populateOptions(channel, "")
        } else {
            populateOptions(channel, UploadData.Channels[selectedIdx])
        }
        channel.val(0);
    }

    function populateOptions(optionNode, options) {
        console.log(optionNode)
        console.log(options)
    }
})();