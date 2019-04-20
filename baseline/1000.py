from random import randint


def action(data, context):
    board = data[1].replace('\n', '')

    while True:
        point = randint(0, 8)
        if board[point] == '0':
            if point % 3 == 0:
                x = 'A'
            elif point % 3 == 1:
                x = 'B'
            else:
                x = 'C'
            y = point // 3 + 1

            return [f"{x}{y}"], {}
