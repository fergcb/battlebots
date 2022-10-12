import sys

test_arena = """..........
..........
..Y.......
..........
..........
.......B..
..........
..........
.......X..
..........
Y hp=10
X hp=10
B x=7 y=5 dir=N
"""


def infer_type(value):
    try:
        value = int(value)
    except Exception:
        pass
    return value


def parse_attrs(attrs):
    pairs = [attr.split("=") for attr in attrs]
    return {entry[0]: infer_type(entry[1]) for entry in pairs}


def parse_entity(line):
    fields = line.split(" ")
    name = fields[0]
    attrs = parse_attrs(fields[1:])
    return {"name": name, **attrs}


def parse_entities(arena):
    source = arena[10*11:]
    return list(map(parse_entity, source.split("\n")[:-1]))


def parse(arena):
    data = {}
    yi = arena.index("Y")
    yx = yi % 11
    yy = yi // 11
    data["me"] = {
        "x": yx,
        "y": yy
    }

    xi = arena.index("X")
    xx = xi % 11
    xy = xi // 11
    data["enemy"] = {
        "x": xx,
        "y": xy
    }

    entities = parse_entities(arena)

    bullets = list(filter(lambda e: e["name"] == "B", entities))
    missiles = list(filter(lambda e: e["name"] == "M", entities))
    landmines = list(filter(lambda e: e["name"] == "L", entities))

    data["bullets"] = bullets
    data["missiles"] = missiles
    data["landmines"] = landmines

    data["me"]["hp"] = entities[0]["hp"]
    data["enemy"]["hp"] = entities[1]["hp"]

    return data


def main(arena):
    data = parse(arena)

    me = data["me"]
    enemy = data["enemy"]

    dx = enemy["x"] - me["x"]
    dy = enemy["y"] - me["y"]

    if abs(dx) > 3 or abs(dy) > 3:
        bx = "W" if dx < 0 else "E" if dx > 0 else ""
        by = "N" if dy < 0 else "S" if dy > 0 else ""

        dir = by + bx
        print(dir)
    else:
        print("P")


if __name__ == "__main__":
    arena = sys.argv[1] if len(sys.argv) > 1 else test_arena
    main(arena)
