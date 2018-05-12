window.addEventListener('load', function () {
    $('[data-toggle="tooltip"]').tooltip();
    var checkbox = $('table tbody input[type="checkbox"]');
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
    checkbox.click(function () {
        if (!this.checked) {
            $("#selectAll").prop("checked", false);
        }
    });

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

        var formData = objectifyForm($form.serializeArray());
        $.ajax({
            url: "/api/meal/" + $form.find('[name="id"]').val(),
            type: 'DELETE',
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
    let id = $(target).closest('tr').find('[data-field="id"]').html();
    let $modal = $('#deleteEmployeeModal');
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
}

