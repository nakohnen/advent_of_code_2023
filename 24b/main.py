import sympy
from dataclasses import dataclass
import argparse

@dataclass
class Hailstone:
    x: int
    y: int
    z: int
    vx: int
    vy: int
    vz: int



def decode_hailstone(s):
    splits = s.split("@")
    left = splits[0].split(",")
    x = left[0].strip()
    y = left[1].strip()
    z = left[2].strip()
    right = splits[1].split(",")
    vx = right[0].strip()
    vy = right[1].strip()
    vz = right[2].strip()
    return int(x), int(y), int(z), int(vx), int(vy), int(vz)

def are_values_unique(my_dict):
    seen_values = set()
    for value in my_dict.values():
        # Check if value is already in the seen set
        if value in seen_values:
            return False  # Duplicate value found
        seen_values.add(value)
    return True  # All values were unique

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Solve hailstone problem")
    parser.add_argument("filename", help="The filename to read the hailstones from")
    args = parser.parse_args()
    filename = args.filename
    hailstones = []
    with open(filename, "r") as f:
        s = f.read()
        for hail in s.split("\n"):
            if hail:
                hailstone = Hailstone(*decode_hailstone(hail))
                hailstones.append(hailstone)

    # for h in hailstones:
    #     print(h)
    print("Hailstones:", len(hailstones))

    x, y, z, vx, vy, vz = sympy.symbols("x y z vx vy vz", integer=True)
    time_variables = []
    for i in range(len(hailstones)):
        time_variables.append(sympy.symbols(f"t{i}", nonnegative=True, real=True, nonzero=True))
    solve_variables = [x, y, z, vx, vy, vz] + time_variables
    print("Solve variables:", solve_variables)
    equations = []
    candidates = []
    print("Creating equations...")
    for h, t in zip(hailstones, time_variables):
        # x + t * vx = hx + t * hvx
        # <=> x - hx = t * (hvx - vx)
        # <=> t = (x - hx) / (hvx - vx)
        eq1 = sympy.Eq(t * (h.vx - vx), x - h.x)
        eq2 = sympy.Eq(t * (h.vy - vy), y - h.y)
        eq3 = sympy.Eq(t * (h.vz - vz), z - h.z)
        equations.append(eq1)
        equations.append(eq2)
        equations.append(eq3)
        cand1 = {x: h.x, vx: h.vx}
        cand2 = {y: h.y, vy: h.vy}
        cand3 = {z: h.z, vz: h.vz}
        candidates.append(cand1)
        candidates.append(cand2)
        candidates.append(cand3)

    # for eq in equations:
    #     print(eq)
    # print()
    print("Equations:", len(equations))
    print("Checking candidates...")
    valid_starts = []
    for c in candidates:
        # print(c)
        valid = True
        for eq in equations:
            subbed = eq.subs(c)
            # print(eq, "=>", subbed)
            if subbed is sympy.false:
                valid = False
                break

        if valid:
            valid_starts.append(c)
    print("Valid starts:", len(valid_starts))
    print()
    
    valid_continued = []
    print("Checking valid starts...")
    for v in valid_starts:
        print("Valid start:", v)
        subbed = [eq.subs(v) for eq in equations]
        subbed2 = [eq for eq in subbed if len(eq.free_symbols) == 1]
        # for s in subbed2:
        #    print("Simplified:", s)
        sol = sympy.solve(subbed2, time_variables)
        # print(sol)
        if are_values_unique(sol):
            new_dict = dict()
            for k, v0 in v.items():
                new_dict[k] = v0
            for k, v0 in sol.items():
                new_dict[k] = v0
            valid_continued.append(new_dict)

    print()
    print("Valid continued:", valid_continued)
    solutions = []
    for v in valid_continued:
        subbed = [eq.subs(v) for eq in equations if eq.subs(v) != sympy.true]
        # for s in subbed:
        #     print(s)
        missing_variables = [var for var in solve_variables if var not in v]
        print("Missing variables:", missing_variables)
        sol = sympy.solve(subbed, missing_variables, dict=True)
        print(f"{sol=}")
        if sol:
            new_dict = dict()
            for k, v0 in v.items():
                new_dict[k] = v0
            for k, v0 in sol.items():
                new_dict[k] = v0
            solutions.append(new_dict)
    print()
    print("Solutions:", len(solutions))
    for s in solutions:
        print(s)
        


