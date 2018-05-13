var cities;
window.addEventListener('load', function () {
    $('#addEmployeeModal [type="submit"]').on('click', function (e) {
        e.preventDefault();
        let $modal = $(this).closest('.modal');
        let $form = $modal.find('form');

        var formData = objectifyForm($form.serializeArray());
        $.ajax({
            url: "/api/meal",
            type: 'POST',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
            },
            data: JSON.stringify(formData),
            dataType: "json",
            success: function (data) {
                $.ajax({
                    url: "/api/meal/all",
                    type: 'GET',
                    beforeSend: function (xhr) {
                        xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                    },
                    success: function (data) {
                        loadTable("#meals-table", data);
                        $form.trigger('reset');
                        $modal.modal('hide');
                    },
                });
            },
            error: function () {
            },
        });
    })

    $('#editEmployeeModal [type="submit"]').on('click', function (e) {
        e.preventDefault();
        let $modal = $(this).closest('.modal');
        let $form = $modal.find('form');

        var formData = objectifyForm($form.serializeArray());
        $.ajax({
            url: "/api/meal/" + $form.find('[name="id"]').val(),
            type: 'POST',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
            },
            data: JSON.stringify(formData),
            dataType: "json",
            success: function (data) {
                $.ajax({
                    url: "/api/meal/all",
                    type: 'GET',
                    beforeSend: function (xhr) {
                        xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                    },
                    success: function (data) {
                        loadTable("#meals-table", data);
                        $form.trigger('reset');
                        $modal.modal('hide');
                    },
                });
            },
            error: function () {
            },
        });
    })

    $('#deleteEmployeeModal [type="submit"]').on('click', function (e) {
        e.preventDefault();
        let $modal = $(this).closest('.modal');
        let $form = $modal.find('form');

        let idsArray = $form.serializeArray();
        $.each(idsArray, function (index, object) {
            $.ajax({
                url: "/api/meal/" + object.value,
                type: 'DELETE',
                beforeSend: function (xhr) {
                    xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
                },
                dataType: "json",
                success: function (data) {
                },
                error: function () {
                },
            });
        })
        $.ajax({
            url: "/api/meal/all",
            type: 'GET',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
            },
            success: function (data) {
                loadTable("#meals-table", data);
                $form.trigger('reset');
                $modal.modal('hide');
            },
        });
    })


    cities = new Bloodhound({
        datumTokenizer: Bloodhound.tokenizers.obj.whitespace('text'),
        queryTokenizer: Bloodhound.tokenizers.whitespace,
       /* prefetch: {
            url: 'static/cities.json',
            cache: false
        }*/
        local: possibleUserTags
    })
    cities.initialize();

    var elt = $('.tagsinput');
    elt.tagsinput({
        itemValue: 'value',
        itemText: 'text',
        typeaheadjs: {
            name: 'cities',
            displayKey: 'text',
            source: cities.ttAdapter()
        }
    });
    elt.tagsinput('add', possibleUserTags[0]);
    elt.tagsinput('add', possibleUserTags[2]);
    elt.tagsinput('add', possibleUserTags[3]);
    elt.tagsinput('add', possibleUserTags[4]);
    elt.tagsinput('add', possibleUserTags[5]);

    $(".twitter-typeahead").css('display', 'inline');

});

function editEntity(target) {
    let id = $(target).closest('tr').find('[data-field="id"]').html();
    let $modal = $('#editEmployeeModal');
    $.ajax({
        url: "/api/meal/" + id,
        type: 'GET',
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
        },
        data: {},
        success: function (data) {
            for (var key in data) {
                if ($modal.find('[name="' + key + '"]').length != 0) {
                    $modal.find('[name="' + key + '"]').val(data[key]);
                }
            }
            $modal.modal('show');
        },
        error: function () {
        },
    });
}

function deleteEntity(target) {
    let $modal = $('#deleteEmployeeModal');
    let $form = $modal.find('form');
    $form.find('input[name="id"]').remove();
    if (target == "checkbox") {
        $('[name="options[]"]:checked').each(function (index, value) {
            let val = $(this).closest('tr').find('[data-field="id"]').html();
            if (val != "") {
                let input = document.createElement("input");
                input.type = "text";
                input.name = "id";
                input.hidden = "hidden";
                input.setAttribute('value', val);
                $form.find('.modal-body')[0].appendChild(input);
            }
        });
    } else {
        let val = $(target).closest('tr').find('[data-field="id"]').html();
        let input = document.createElement("input");
        input.type = "text";
        input.name = "id";
        input.hidden = "hidden";
        input.setAttribute('value', val);
        $form.find('.modal-body')[0].appendChild(input);
        $modal.modal('show');
    }
}

function getCookieValue(a) {
    var b = document.cookie.match('(^|;)\\s*' + a + '\\s*=\\s*([^;]+)');
    return b ? b.pop() : '';
}

function objectifyForm(formArray) {//serialize data function

    var returnArray = {};
    for (var i = 0; i < formArray.length; i++) {
        returnArray[formArray[i]['name']] = formArray[i]['value'];
    }
    return returnArray;
}

function loadTable(table, data) {
    $("tbody tr:not(:first)").remove();
    let $table = $(table);
    let dataArray = data;
    let $firstRow = $table.find('tbody tr:first');

    $.each(dataArray, function (index, object) {
        let $row = $firstRow.clone();
        $table.find('tbody tr:last').after($row);
        $.each(object, function (key, value) {
            let $td = $row.find('[data-field="' + key + '"]');
            if ($td.length != 0) {
                $td.html(value);
            }
            //alert( key + ": " + value );
        });
        $row.show();
    });

    $('[data-toggle="tooltip"]').tooltip();

    $("#selectAll").unbind('click');
    $("#selectAll").click(function () {
        if (this.checked) {
            checkbox.each(function () {
                this.checked = true;
            });
        } else {
            checkbox.each(function () {
                this.checked = false;
            });
        }
    });
    let checkboxSpan = $('table tbody .custom-checkbox').on('click', function (e) {
        if (e.target.type == "checkbox") {
            return
        }
        let chk = $(this).find('input[type="checkbox"]');
        if (chk.prop('checked') == true) {
            chk.prop("checked", false);
        } else {
            chk.prop("checked", true);
        }
    });
    let checkbox = $('table tbody input[type="checkbox"]');
    checkbox.removeAttr('checked');
    checkbox.click(function () {
        if (!this.checked) {
            $("#selectAll").prop("checked", false);
        }
    });
}

