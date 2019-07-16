#include <cstdio>
#include <map>
#include <string>
#include <vector>

class potsherd_type {
public:
    std::string name;
    enum type_category {
        NUMBER = 0,
        STRING = 1,
        FIXARRAY = 2,
        DYNARRAY = 3,
        STRUCT = 4,
        UNION = 5,
        ENUM = 6
    } category;

    potsherd_type() { }
    potsherd_type(std::string name, type_category category)
        : name(name), category(category) { }
};

class potsherd_type_fixarray : public potsherd_type {
public:
    potsherd_type *elm_type;
    std::vector<int> dims;
};

class potsherd_type_dynarray : public potsherd_type {
public:
    potsherd_type *elm_type;
};

class potsherd_type_composite : public potsherd_type {
public:
    std::vector<potsherd_type *> members;
};

class potsherd_type_enum : public potsherd_type {
public:
    struct entry {
        std::string name;
        int value;
    };
    std::vector<entry> entries;
};

static inline bool is_identifier(char ch)
{
    return (ch >= '0' && ch <= '9') ||
        (ch >= 'A' && ch <= 'Z') ||
        (ch >= 'a' && ch <= 'z') || ch == '_';
}

static inline bool is_token(char ch)
{
    return is_identifier(ch) || ch == '[' || ch == ']';
}

static inline bool is_keyword(const std::string &s)
{
    return s == "struct" || s == "union" || s == "enum";
}

struct parse_state {
    const char *s;
    int pos;
    int last;
    int err;

    std::map<std::string, potsherd_type *> map;

    parse_state()
    {
        // Set up built-in types
        map["int8"] = new potsherd_type("int8", potsherd_type::NUMBER);
        map["bool"] = new potsherd_type("bool", potsherd_type::NUMBER);
        map["string"] = new potsherd_type("string", potsherd_type::STRING);
    }

    inline std::string next_token()
    {
        if (err) return std::string();

        while (s[pos] != '\0' && !is_token(s[pos])) pos++;
        last = pos;
        if (s[pos] == '\0') return std::string();

        while (is_token(s[pos])) pos++;
        return std::string(s + last, pos - last);
    }

    void back_token()
    {
        pos = last;
    }

    potsherd_type *parse_type(std::string s)
    {
        int p = 0;
        while (p < s.length() && is_identifier(s[p])) ++p;
        std::string atomic_type(s, 0, p);

        std::map<std::string, potsherd_type *>::iterator it;
        if ((it = map.find(atomic_type)) == map.end()) return NULL;

        printf("Parsing type: [%s]\n", s.c_str());
        printf("Atomic type: [%s]\n", atomic_type.c_str());

        return (potsherd_type *)0x1;    // = =
    }

    inline void parse_enum()
    {
        if (err) return;

        potsherd_type_enum *ret = new potsherd_type_enum();
        std::string name = next_token();
        ret->name = name;
        ret->category = potsherd_type::ENUM;

        std::string token;
        while (!(token = next_token()).empty()) {
            if (is_keyword(token)) break;
            if (parse_type(next_token()) != NULL) {
                err = 1;
                free(ret);
                return;
            }
        }

        back_token();
        map[name] = ret;
    }

    inline void parse_composite(bool is_union)
    {
        if (err) return;

        potsherd_type_composite *ret = new potsherd_type_composite();
        std::string name = next_token();
        ret->name = name;
        ret->category = (is_union ? potsherd_type::UNION : potsherd_type::STRUCT);

        potsherd_type *type;
        while ((type = parse_type(next_token())) != NULL) {
            std::string member_name = next_token();
        }

        back_token();
        map[name] = ret;
    }

    static std::vector<potsherd_type *> parse(const char *s)
    {
        parse_state state;
        state.s = s;
        state.pos = 0;
        state.err = 0;

        std::string token;
        while (!(token = state.next_token()).empty()) {
            puts(token.c_str());
            if (token == "struct" || token == "union") {
                state.parse_composite(token == "union");
            } else if (token == "enum") {
                state.parse_enum();
            }
        }

        std::vector<potsherd_type *> ret;
        for (std::map<std::string, potsherd_type *>::iterator
            it = state.map.begin();
            it != state.map.end(); it++)
        {
            ret.push_back(it->second);
        }
        return ret;
    }
};

int main()
{
    FILE *f = fopen("potsherd.txt", "r");
    if (!f) {
        puts("Cannot open file > <");
        return 1;
    }

    fseek(f, 0, SEEK_END);
    long size = ftell(f);

    char *buf = new char[size + 1];
    fseek(f, 0, SEEK_SET);
    if (fread(buf, size, 1, f) < 1) {
        puts("Cannot read from file > <");
        return 2;
    }
    buf[size] = '\0';

    std::vector<potsherd_type *> list = parse_state::parse(buf);
    for (int i = 0; i < list.size(); i++)
        puts(list[i]->name.c_str());

    fclose(f);
    return 0;
}
